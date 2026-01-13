package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/middleware"
)

func TestLogger(t *testing.T) {
	r := router.New()

	r.GET("/test", middleware.Logger(func(ctx router.Context) error {
		return ctx.String(http.StatusOK, "OK")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %q", w.Body.String())
	}
}

func TestRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(router.Context) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "handles panic",
			handler: func(ctx router.Context) error {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
		{
			name: "normal execution",
			handler: func(ctx router.Context) error {
				return ctx.String(http.StatusOK, "OK")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.New()

			r.GET("/test", middleware.Recovery(tt.handler))

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	r := router.New()

	r.GET("/test", middleware.SecurityHeaders(func(ctx router.Context) error {
		return ctx.String(http.StatusOK, "OK")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	headers := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, expected := range headers {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("header %s: expected %q, got %q", header, expected, got)
		}
	}
}

func TestRequestID(t *testing.T) {
	r := router.New()

	var capturedID string
	r.GET("/test", middleware.RequestID(func(ctx router.Context) error {
		capturedID = ctx.Get("request_id").(string)
		return ctx.String(http.StatusOK, "OK")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	responseID := w.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("expected X-Request-ID header to be set")
	}

	if capturedID == "" {
		t.Error("expected request_id to be stored in context")
	}

	if responseID != capturedID {
		t.Errorf("response header and context ID don't match: %q != %q", responseID, capturedID)
	}
}

func TestCORS(t *testing.T) {
	allowedOrigins := []string{"https://example.com"}
	corsMiddleware := middleware.CORS(allowedOrigins)

	t.Run("allowed origin", func(t *testing.T) {
		r := router.New()
		r.GET("/test", corsMiddleware(func(ctx router.Context) error {
			return ctx.String(http.StatusOK, "OK")
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
			t.Errorf("expected CORS header for allowed origin, got %q", got)
		}
	})

	t.Run("disallowed origin", func(t *testing.T) {
		r := router.New()
		r.GET("/test", corsMiddleware(func(ctx router.Context) error {
			return ctx.String(http.StatusOK, "OK")
		}))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "https://evil.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Errorf("expected no CORS header for disallowed origin, got %q", got)
		}
	})

	t.Run("preflight request", func(t *testing.T) {
		r := router.New()
		r.OPTIONS("/test", corsMiddleware(func(ctx router.Context) error {
			return ctx.String(http.StatusOK, "")
		}))

		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for preflight, got %d", w.Code)
		}
	})
}
