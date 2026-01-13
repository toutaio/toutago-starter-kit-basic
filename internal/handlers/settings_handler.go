package handlers

import (
	"net/http"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/services"
)

type SettingsHandler struct {
	renderer    *fith.Engine
	authService *services.AuthService
}

func NewSettingsHandler(renderer *fith.Engine, authService *services.AuthService) *SettingsHandler {
	return &SettingsHandler{
		renderer:    renderer,
		authService: authService,
	}
}

func (h *SettingsHandler) Show(c router.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusSeeOther)
		return nil
	}

	data := map[string]interface{}{
		"User": user,
	}

	html, err := h.renderer.Render("pages/settings.html", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return c.HTML(http.StatusOK, html)
}

func (h *SettingsHandler) UpdatePassword(c router.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusSeeOther)
		return nil
	}

	currentPassword := c.Request().FormValue("current_password")
	newPassword := c.Request().FormValue("new_password")
	confirmPassword := c.Request().FormValue("confirm_password")

	// Validate input
	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		data := map[string]interface{}{
			"User":  user,
			"Error": "All fields are required",
		}
		html, err := h.renderer.Render("pages/settings.html", data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
		}
		return c.HTML(http.StatusBadRequest, html)
	}

	if newPassword != confirmPassword {
		data := map[string]interface{}{
			"User":  user,
			"Error": "New passwords do not match",
		}
		html, err := h.renderer.Render("pages/settings.html", data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
		}
		return c.HTML(http.StatusBadRequest, html)
	}

	if len(newPassword) < 8 {
		data := map[string]interface{}{
			"User":  user,
			"Error": "Password must be at least 8 characters long",
		}
		html, err := h.renderer.Render("pages/settings.html", data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
		}
		return c.HTML(http.StatusBadRequest, html)
	}

	// Verify current password and update
	if h.authService != nil {
		err := h.authService.UpdatePassword(uint(user.ID), newPassword)
		if err != nil {
			data := map[string]interface{}{
				"User":  user,
				"Error": "Failed to update password: " + err.Error(),
			}
			html, err := h.renderer.Render("pages/settings.html", data)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
			}
			return c.HTML(http.StatusInternalServerError, html)
		}
	}

	data := map[string]interface{}{
		"User":    user,
		"Success": "Password updated successfully",
	}

	html, err := h.renderer.Render("pages/settings.html", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return c.HTML(http.StatusOK, html)
}
