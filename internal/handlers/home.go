package handlers

import (
	"net/http"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
)

// HomeHandler handles home page requests.
type HomeHandler struct {
	renderer *fith.Engine
}

// NewHomeHandler creates a new home handler.
func NewHomeHandler(renderer *fith.Engine) *HomeHandler {
	return &HomeHandler{renderer: renderer}
}

// Index handles the home page.
func (h *HomeHandler) Index(ctx router.Context) error {
	data := map[string]interface{}{
		"title": "Home",
	}

	html, err := h.renderer.Render("pages/home.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}
