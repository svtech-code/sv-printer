package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultHost           = "127.0.0.1"
	DefaultPort           = 9876
	DefaultMaxPayloadSize = int64(5 * 1024 * 1024)
)

var DefaultAllowedOrigin = ""

type ManualPrinter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Protocol   string `json:"protocol"`
	Connection string `json:"connection,omitempty"`
	BaudRate   int    `json:"baud_rate,omitempty"`
}

type Config struct {
	Host           string          `json:"host"`
	Port           int             `json:"port"`
	Token          string          `json:"token"`
	AllowedOrigins []string        `json:"allowed_origins"`
	MaxPayloadSize int64           `json:"max_payload_size"`
	Printers       []ManualPrinter `json:"printers"`
	TokenGenerated bool            `json:"-"`
	Path           string          `json:"-"`
	LicensePath    string          `json:"-"`
	Headless       bool            `json:"headless"`
}

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }

func (s *stringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func Load(args []string) (Config, error) {
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)

	configPath := fs.String("config", envStr("SV_PRINT_CONFIG", ""), "path to the config file")
	licensePath := fs.String("license", envStr("SV_PRINT_LICENSE", ""), "path to the license file")
	host := fs.String("host", "", "address to bind the local API")
	port := fs.Int("port", 0, "port to bind the local API")
	token := fs.String("token", "", "auth token for the local API")
	maxPayload := fs.Int64("max-payload", 0, "maximum payload size in bytes")
	headless := fs.Bool("headless", false, "run without system tray icon")

	var origins stringList
	fs.Var(&origins, "origin", "allowed CORS origin (repeatable)")

	var printers stringList
	fs.Var(&printers, "printer", "manual network printer as name@address (repeatable)")

	var serialPrinters stringList
	fs.Var(&serialPrinters, "serial", "manual serial printer as name@port[@baud] (repeatable)")

	var systemPrinters stringList
	fs.Var(&systemPrinters, "system", "manual system printer as name@printer_id (repeatable)")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	var cfg Config
	if *configPath != "" {
		if err := loadFile(*configPath, &cfg); err != nil && !os.IsNotExist(err) {
			return Config{}, err
		}
	}

	if *host != "" {
		cfg.Host = *host
	} else if env := os.Getenv("SV_PRINT_HOST"); env != "" {
		cfg.Host = env
	} else if cfg.Host == "" {
		cfg.Host = DefaultHost
	}

	if *port != 0 {
		cfg.Port = *port
	} else if env := os.Getenv("SV_PRINT_PORT"); env != "" {
		if n, err := strconv.Atoi(env); err == nil {
			cfg.Port = n
		}
	}
	if cfg.Port == 0 {
		cfg.Port = DefaultPort
	}

	if *token != "" {
		cfg.Token = *token
	} else if env := os.Getenv("SV_PRINT_TOKEN"); env != "" {
		cfg.Token = env
	}

	if *headless {
		cfg.Headless = true
	}

	if *maxPayload != 0 {
		cfg.MaxPayloadSize = *maxPayload
	} else if env := os.Getenv("SV_PRINT_MAX_PAYLOAD"); env != "" {
		if n, err := strconv.ParseInt(env, 10, 64); err == nil {
			cfg.MaxPayloadSize = n
		}
	}
	if cfg.MaxPayloadSize == 0 {
		cfg.MaxPayloadSize = DefaultMaxPayloadSize
	}

	if len(origins) > 0 {
		cfg.AllowedOrigins = []string(origins)
	} else if len(cfg.AllowedOrigins) == 0 {
		if env := os.Getenv("SV_PRINT_ALLOWED_ORIGINS"); env != "" {
			cfg.AllowedOrigins = strings.Split(env, ",")
		} else if DefaultAllowedOrigin != "" {
			cfg.AllowedOrigins = []string{DefaultAllowedOrigin}
		}
	}

	if len(printers) > 0 {
		cfg.Printers = nil
		for _, p := range printers {
			mp, err := parsePrinter(p)
			if err != nil {
				return Config{}, err
			}
			cfg.Printers = append(cfg.Printers, mp)
		}
	}

	for _, p := range serialPrinters {
		mp, err := parseSerialPrinter(p)
		if err != nil {
			return Config{}, err
		}
		cfg.Printers = append(cfg.Printers, mp)
	}

	for _, p := range systemPrinters {
		mp, err := parseSystemPrinter(p)
		if err != nil {
			return Config{}, err
		}
		cfg.Printers = append(cfg.Printers, mp)
	}

	if cfg.Token == "" {
		generated, err := generateToken()
		if err != nil {
			return Config{}, err
		}
		cfg.Token = generated
		cfg.TokenGenerated = true
	}

	cfg.Path = *configPath
	if cfg.Path == "" {
		cfg.Path = DefaultPath()
	}

	cfg.LicensePath = *licensePath
	if cfg.LicensePath == "" {
		cfg.LicensePath = filepath.Join(filepath.Dir(cfg.Path), "license.key")
	}

	return cfg, nil
}

func Save(cfg Config) error {
	path := cfg.Path
	if path == "" {
		path = DefaultPath()
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}

func DefaultPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "sv-printer-config.json"
	}
	return filepath.Join(dir, "sv-printer", "config.json")
}

func (c Config) LogPath() string {
	return filepath.Join(filepath.Dir(c.Path), "sv-printer.log")
}

func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cfg)
}

func parsePrinter(s string) (ManualPrinter, error) {
	name := s
	address := s
	if i := strings.LastIndex(s, "@"); i >= 0 {
		name = s[:i]
		address = s[i+1:]
	}
	if address == "" {
		return ManualPrinter{}, fmt.Errorf("invalid printer %q: missing address", s)
	}
	return ManualPrinter{
		ID:         "net-" + address,
		Name:       name,
		Address:    address,
		Protocol:   "escpos",
		Connection: "network",
	}, nil
}

func parseSerialPrinter(s string) (ManualPrinter, error) {
	parts := strings.Split(s, "@")
	var name, port string
	baud := 9600

	switch len(parts) {
	case 1:
		port = parts[0]
	case 2:
		name = parts[0]
		port = parts[1]
	case 3:
		name = parts[0]
		port = parts[1]
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			return ManualPrinter{}, fmt.Errorf("invalid baud rate %q: %v", parts[2], err)
		}
		baud = n
	default:
		return ManualPrinter{}, fmt.Errorf("invalid serial printer %q", s)
	}

	if port == "" {
		return ManualPrinter{}, fmt.Errorf("invalid serial printer %q: missing port", s)
	}

	return ManualPrinter{
		ID:         "serial-" + port,
		Name:       name,
		Address:    port,
		Protocol:   "escpos",
		Connection: "serial",
		BaudRate:   baud,
	}, nil
}

func parseSystemPrinter(s string) (ManualPrinter, error) {
	name := s
	address := s
	if i := strings.LastIndex(s, "@"); i >= 0 {
		name = s[:i]
		address = s[i+1:]
	}
	if address == "" {
		return ManualPrinter{}, fmt.Errorf("invalid system printer %q: missing printer_name", s)
	}
	return ManualPrinter{
		ID:         "system-" + address,
		Name:       name,
		Address:    address,
		Protocol:   "escpos",
		Connection: "system",
	}, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
