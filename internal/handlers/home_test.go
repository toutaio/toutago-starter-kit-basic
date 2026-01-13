package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/handlers"
)

func TestHomeHandler(t *testing.T) {
	// Create renderer
	config := &fith.Config{
		TemplateDir: "../../templates",
	}
	renderer, err := fith.New(config)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	r := router.New()
	handler := handlers.NewHomeHandler(renderer)
	r.GET("/", handler.Index)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	// Check that we got some HTML response
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected HTML content type, got %s", contentType)
	}

	if w.Body.Len() == 0 {
		t.Error("expected response body, got empty")
	}
}
