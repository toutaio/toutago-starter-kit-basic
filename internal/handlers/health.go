// Package handlers provides HTTP request handlers.
package handlers

import (
	"database/sql"
	"net/http"

	router "github.com/toutaio/toutago-cosan-router"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
}

// Check handles the health check endpoint.
func (h *HealthHandler) Check(ctx router.Context) error {
	response := HealthResponse{
		Status: "healthy",
	}

	// Check database if available
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			response.Status = "unhealthy"
			response.Database = "disconnected"
			return ctx.JSON(http.StatusServiceUnavailable, response)
		}
		response.Database = "connected"
	} else {
		response.Database = "not configured"
	}

	return ctx.JSON(http.StatusOK, response)
}
