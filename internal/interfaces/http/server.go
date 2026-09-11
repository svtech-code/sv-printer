package http

import (
	"fmt"
	"net/http"
)

type ServerConfig struct {
	Host           string
	Port           int
	AllowedOrigins []string
	Token          string
	MaxPayloadSize int64
}

func NewServer(config ServerConfig, api *API) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.HealthHandler)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/v1/info", api.InfoHandler)
	protectedMux.HandleFunc("GET /api/v1/printers", api.PrintersHandler)
	protectedMux.HandleFunc("POST /api/v1/printers/discover", api.DiscoverHandler)
	protectedMux.HandleFunc("GET /api/v1/printers/{id}", api.GetPrinterHandler)
	protectedMux.HandleFunc("POST /api/v1/printers/{id}/test", api.TestPrinterHandler)
	protectedMux.HandleFunc("POST /api/v1/print", api.PrintHandler)
	protectedMux.HandleFunc("POST /api/v1/print/receipt", api.PrintReceiptHandler)
	protectedMux.HandleFunc("GET /api/v1/jobs/{id}", api.GetJobHandler)
	protectedMux.HandleFunc("GET /api/v1/events", api.EventsHandler)

	protectedHandler := WithAuth(config.Token, protectedMux)
	protectedHandler = WithPayloadLimit(config.MaxPayloadSize, protectedHandler)
	protectedHandler = WithCORS(config.AllowedOrigins, protectedHandler)

	mux.Handle("/api/", protectedHandler)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
