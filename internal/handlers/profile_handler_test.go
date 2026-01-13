package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
)

func TestProfileHandler_Show(t *testing.T) {
	router := cosan.New()
	handler := handlers.NewProfileHandler(nil, nil)
	
	router.GET("/profile", handler.Show)

	t.Run("requires authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusUnauthorized {
			t.Errorf("expected redirect or unauthorized, got %d", w.Code)
		}
	})

	// TODO: Add test for authenticated user once we have proper context creation
}

func TestProfileHandler_Update(t *testing.T) {
	router := cosan.New()
	handler := handlers.NewProfileHandler(nil, nil)
	
	router.POST("/profile", handler.Update)

	t.Run("requires authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/profile", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusUnauthorized {
			t.Errorf("expected redirect or unauthorized, got %d", w.Code)
		}
	})
}
