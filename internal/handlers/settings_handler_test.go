package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
)

func TestSettingsHandler_Show(t *testing.T) {
	router := cosan.New()
	handler := handlers.NewSettingsHandler(nil, nil)
	
	router.GET("/settings", handler.Show)

	t.Run("requires authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/settings", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusUnauthorized {
			t.Errorf("expected redirect or unauthorized, got %d", w.Code)
		}
	})

	// TODO: Add test for authenticated user once we have proper context creation
}

func TestSettingsHandler_UpdatePassword(t *testing.T) {
	router := cosan.New()
	handler := handlers.NewSettingsHandler(nil, nil)
	
	router.POST("/settings/password", handler.UpdatePassword)

	t.Run("requires authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/settings/password", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusUnauthorized {
			t.Errorf("expected redirect or unauthorized, got %d", w.Code)
		}
	})
}
