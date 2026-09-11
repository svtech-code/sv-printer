package http

import (
	"crypto/subtle"
	"net/http"
	"strings"

	domainErrors "sv-print/internal/domain/errors"
)

func WithAuth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqToken string

		if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			reqToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else if qp := r.URL.Query().Get("token"); qp != "" {
			reqToken = qp
		}

		if reqToken == "" || subtle.ConstantTimeCompare([]byte(reqToken), []byte(token)) != 1 {
			s, c, m := mapDomainError(domainErrors.ErrUnauthorized)
			writeError(w, s, c, m)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithCORS(allowedOrigins []string, next http.Handler) http.Handler {
	originMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		originMap[o] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				s, c, m := mapDomainError(domainErrors.ErrOriginNotAllowed)
				writeError(w, s, c, m)
				return
			}
		}

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithPayloadLimit(limit int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		next.ServeHTTP(w, r)
	})
}
