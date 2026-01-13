package handlers

import (
	"net/http"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repositories"
)

type ProfileHandler struct {
	renderer *fith.Engine
	userRepo repositories.UserRepository
}

func NewProfileHandler(renderer *fith.Engine, userRepo repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{
		renderer: renderer,
		userRepo: userRepo,
	}
}

func (h *ProfileHandler) Show(c router.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusSeeOther)
		return nil
	}

	data := map[string]interface{}{
		"User": user,
	}

	html, err := h.renderer.Render("pages/profile.html", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return c.HTML(http.StatusOK, html)
}

func (h *ProfileHandler) Update(c router.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusSeeOther)
		return nil
	}

	// Get form data
	email := c.Request().FormValue("email")
	
	if email == "" {
		data := map[string]interface{}{
			"User":  user,
			"Error": "Email is required",
		}
		html, err := h.renderer.Render("pages/profile.html", data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
		}
		return c.HTML(http.StatusBadRequest, html)
	}

	// Update user
	user.Email = email
	
	if h.userRepo != nil {
		err := h.userRepo.Update(user)
		if err != nil {
			data := map[string]interface{}{
				"User":  user,
				"Error": "Failed to update profile: " + err.Error(),
			}
			html, err := h.renderer.Render("pages/profile.html", data)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
			}
			return c.HTML(http.StatusInternalServerError, html)
		}
	}

	data := map[string]interface{}{
		"User":    user,
		"Success": "Profile updated successfully",
	}

	html, err := h.renderer.Render("pages/profile.html", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return c.HTML(http.StatusOK, html)
}
