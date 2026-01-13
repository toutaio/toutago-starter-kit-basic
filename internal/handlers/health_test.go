package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		db             *sql.DB
		expectedStatus int
		checkBody      bool
	}{
		{
			name:           "healthy with nil db (not yet connected)",
			db:             nil,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.New()
			handler := handlers.NewHealthHandler(tt.db)
			r.GET("/health", handler.Check)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkBody && w.Body.Len() == 0 {
				t.Error("expected response body, got empty")
			}
		})
	}
}
