package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"sv-printer/internal/application/discovery"
	"sv-printer/internal/application/events"
	"sv-printer/internal/application/printing"
	"sv-printer/internal/config"
	"sv-printer/internal/deviceid"
	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/domain/printer"
	"sv-printer/internal/infrastructure/network"
	"sv-printer/internal/infrastructure/serial"
	"sv-printer/internal/infrastructure/system"
	"sv-printer/internal/infrastructure/transport"
	"sv-printer/internal/infrastructure/usb"
	apphttp "sv-printer/internal/interfaces/http"
	"sv-printer/internal/licensing"
	"sv-printer/internal/receipt"
)

var Version = "0.1.0"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	args := os.Args[1:]

	var err error
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		err = runAgent(ctx, args, os.Stdout)
	} else {
		err = runCLI(ctx, args, os.Stdout)
	}

	if err != nil {
		slog.Error("application error", "error", err.Error())
		os.Exit(1)
	}
}

func runCLI(ctx context.Context, args []string, stdout io.Writer) error {
	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "SV Print v%s\n", Version)
		return nil
	case "status":
		fmt.Fprintln(stdout, "SV Print status: OK")
		return nil
	case "printers":
		cfg, err := config.Load(args[1:])
		if err != nil {
			return err
		}
		if len(cfg.Printers) == 0 {
			fmt.Fprintln(stdout, "No printers configured yet.")
			return nil
		}
		fmt.Fprintln(stdout, "SV Print")
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "ID\tNAME\tADDRESS")
		for _, p := range cfg.Printers {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", p.ID, p.Name, p.Address)
		}
		return nil
	case "config":
		cfg, err := config.Load(args[1:])
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "config file: %s\n", cfg.Path)
		fmt.Fprintf(stdout, "host: %s\n", cfg.Host)
		fmt.Fprintf(stdout, "port: %d\n", cfg.Port)
		fmt.Fprintf(stdout, "token: %s\n", redactToken(cfg.Token))
		fmt.Fprintf(stdout, "max_payload: %d\n", cfg.MaxPayloadSize)
		if len(cfg.AllowedOrigins) > 0 {
			fmt.Fprintf(stdout, "allowed_origins: %s\n", strings.Join(cfg.AllowedOrigins, ","))
		}
		if len(cfg.Printers) > 0 {
			for _, p := range cfg.Printers {
				fmt.Fprintf(stdout, "printer: %s (%s) %s\n", p.ID, p.Name, p.Address)
			}
		} else {
			fmt.Fprintln(stdout, "printers: none")
		}
		return nil
	case "discover":
		cfg, err := config.Load(args[1:])
		if err != nil {
			return err
		}
		printers, err := buildRegistry(cfg).DiscoverAll(ctx)
		if err != nil {
			return err
		}
		if len(printers) == 0 {
			fmt.Fprintln(stdout, "No printers discovered.")
			return nil
		}
		fmt.Fprintln(stdout, "ID\tNAME\tCONNECTION")
		for _, p := range printers {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", p.ID, p.Name, p.Connection)
		}
		return nil
	case "test":
		if len(args) < 2 {
			return fmt.Errorf("usage: sv-printer test <printer-id>")
		}
		cfg, err := config.Load(nil)
		if err != nil {
			return err
		}
		wm, err := checkPrintAllowed(licensing.FromFile(cfg.LicensePath))
		if err != nil {
			return err
		}
		doc := receipt.Document{
			Cut: true,
			Lines: []receipt.Line{
				{Text: "SV Print", Style: receipt.Style{Bold: true, Align: "center"}},
				{Text: "Test receipt"},
			},
		}
		if err := sendToPrinter(ctx, cfg, args[1], doc.Build(wm).Bytes()); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Test receipt sent to %s\n", args[1])
		return nil
	case "print":
		if len(args) < 3 {
			return fmt.Errorf("usage: sv-printer print <printer-id> <receipt.json>")
		}
		cfg, err := config.Load(nil)
		if err != nil {
			return err
		}
		wm, err := checkPrintAllowed(licensing.FromFile(cfg.LicensePath))
		if err != nil {
			return err
		}
		data, err := os.ReadFile(args[2])
		if err != nil {
			return err
		}
		var doc receipt.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("invalid receipt file: %v", err)
		}
		if err := doc.Validate(); err != nil {
			return err
		}
		if err := sendToPrinter(ctx, cfg, args[1], doc.Build(wm).Bytes()); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Receipt sent to %s\n", args[1])
		return nil
	case "logs":
		cfg, err := config.Load(args[1:])
		if err != nil {
			return err
		}
		if err := tailLogs(stdout, cfg.LogPath()); err != nil {
			return err
		}
		return nil
	case "doctor":
		cfg, err := config.Load(args[1:])
		if err != nil {
			return err
		}
		runDoctor(ctx, stdout, cfg)
		return nil
	case "device-id":
		id, err := deviceid.ID()
		if err != nil {
			return err
		}
		fmt.Fprintln(stdout, id)
		return nil
	case "service":
		return fmt.Errorf("service is not implemented yet")
	default:
		return fmt.Errorf("unknown subcommand: %s", args[0])
	}
}

func redactToken(t string) string {
	if t == "" {
		return "(none)"
	}
	return "********"
}

func runAgent(ctx context.Context, args []string, stdout io.Writer) error {
	cfg, err := config.Load(args)
	if err != nil {
		return err
	}

	setupLogFile(cfg)

	queue := printing.NewInMemoryQueue()
	registry := buildRegistry(cfg)

	bus := events.NewBus()

	licState := licensing.FromFile(cfg.LicensePath)

	api := apphttp.NewAPI(queue, registry, Version, bus).WithLicense(licState)

	server := apphttp.NewServer(apphttp.ServerConfig{
		Host:           cfg.Host,
		Port:           cfg.Port,
		AllowedOrigins: cfg.AllowedOrigins,
		Token:          cfg.Token,
		MaxPayloadSize: cfg.MaxPayloadSize,
	}, api)

	worker := printing.NewWorker(queue, transportFactory(cfg), bus)

	if cfg.TokenGenerated {
		fmt.Fprintf(stdout, "SV Print initialized.\n\nToken:\n%s\n\n", cfg.Token)
		if err := config.Save(cfg); err != nil {
			slog.Warn("failed to persist config", "error", err.Error())
		}
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go worker.Start(ctx)

	slog.Info("print worker started")
	bus.Publish(events.NewEvent(events.AgentStatus, map[string]any{"status": "started"}))

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	bus.Publish(events.NewEvent(events.AgentStatus, map[string]any{"status": "stopped"}))
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	worker.Stop()

	return nil
}

func registerManualPrinters(registry *discovery.Registry, cfg config.Config) {
	var printers []printer.Printer
	for _, p := range cfg.Printers {
		conn := printer.ConnectionNetwork
		if p.Connection == "serial" {
			conn = printer.ConnectionSerial
		} else if p.Connection == "system" {
			conn = printer.ConnectionSystem
		}
		printers = append(printers, printer.Printer{
			ID:         p.ID,
			Name:       p.Name,
			Connection: conn,
			Address:    p.Address,
			Protocol:   printer.ProtocolEscpos,
			Status:     printer.StatusUnknown,
		})
	}
	if len(printers) > 0 {
		registry.Register(discovery.NewStaticDiscoverer(printers))
	}
}

func buildRegistry(cfg config.Config) *discovery.Registry {
	registry := discovery.NewRegistry()
	registerManualPrinters(registry, cfg)
	// Hardware discoverers
	registry.Register(usb.NewDiscoverer())
	registry.Register(serial.NewDiscoverer())
	registry.Register(system.NewDiscoverer())
	return registry
}

func transportFactory(cfg config.Config) printing.TransportFactory {
	return func(printerID string) (transport.PrinterTransport, error) {
		for _, p := range cfg.Printers {
			if p.ID == printerID {
				if p.Connection == "serial" {
					return serial.NewSerialTransport(p.Address, p.BaudRate), nil
				}
				if p.Connection == "system" {
					return system.NewSystemTransport(p.Address), nil
				}
				return network.NewTCPTransport(p.Address, 5*time.Second), nil
			}
		}
		return nil, domainErrors.ErrPrinterNotFound
	}
}

func sendToPrinter(ctx context.Context, cfg config.Config, printerID string, payload []byte) error {
	t, err := transportFactory(cfg)(printerID)
	if err != nil {
		return err
	}
	defer t.Close()

	if err := t.Open(ctx); err != nil {
		return err
	}
	return t.Write(ctx, payload)
}

func checkPrintAllowed(state *licensing.State) (string, error) {
	if !state.AllowPrint() {
		return "", domainErrors.ErrLicenseQuotaExceeded
	}
	return state.Watermark(), nil
}

func setupLogFile(cfg config.Config) {
	f, err := os.OpenFile(cfg.LogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stderr, f), nil)))
}

func tailLogs(stdout io.Writer, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stdout, "No logs found at %s\n", path)
		return nil
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	const n = 100
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	for _, l := range lines {
		fmt.Fprintln(stdout, l)
	}
	return nil
}

func runDoctor(ctx context.Context, stdout io.Writer, cfg config.Config) {
	fmt.Fprintln(stdout, "SV Print Doctor")
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "✓ Operating system: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(stdout, "✓ Config file: %s\n", cfg.Path)
	fmt.Fprintln(stdout, "✓ Config valid")

	printers, _ := buildRegistry(cfg).DiscoverAll(ctx)
	if len(printers) == 0 {
		fmt.Fprintln(stdout, "✗ No printers detected")
		return
	}
	fmt.Fprintf(stdout, "✓ Printers detected: %d\n", len(printers))
	for _, p := range printers {
		fmt.Fprintf(stdout, "  ✓ %s (%s)\n", p.Name, p.ID)
	}
}
