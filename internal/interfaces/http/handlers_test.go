package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sv-printer/internal/application/discovery"
	"sv-printer/internal/application/printing"
)

func TestHealthHandler(t *testing.T) {
	q := printing.NewInMemoryQueue()
	r := discovery.NewRegistry()
	api := NewAPI(q, r, "1.0.0")

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.HealthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)

	if response["status"] != "ok" {
		t.Errorf("expected status ok, got %v", response["status"])
	}
}

func TestWithAuth(t *testing.T) {
	handler := WithAuth("secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("missing auth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %v", rr.Code)
		}
	})

	t.Run("wrong header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer wrong")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %v", rr.Code)
		}
	})

	t.Run("correct header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer secret")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %v", rr.Code)
		}
	})

	t.Run("correct query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?token=secret", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %v", rr.Code)
		}
	})

	t.Run("wrong query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?token=wrong", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %v", rr.Code)
		}
	})

	t.Run("header takes precedence over query", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/?token=wrong", nil)
		req.Header.Set("Authorization", "Bearer secret")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %v", rr.Code)
		}
	})
}
