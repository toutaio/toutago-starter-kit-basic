package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
	"github.com/toutaio/toutago-starter-kit-basic/internal/helpers"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/service"
)

// PageHandler handles page-related requests
type PageHandler struct {
	pageService *service.PageService
	renderer    *fith.Engine
}

// NewPageHandler creates a new page handler
func NewPageHandler(pageService *service.PageService, renderer *fith.Engine) *PageHandler {
	return &PageHandler{
		pageService: pageService,
		renderer:    renderer,
	}
}

// Index displays list of pages
func (h *PageHandler) Index(ctx router.Context) error {
	// Parse pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20
	offset := (page - 1) * perPage

	// Get pages
	pages, err := h.pageService.ListPublishedPages(ctx.Request().Context(), perPage, offset)
	if err != nil {
		log.Printf("Error listing pages: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error loading pages")
	}

	data := map[string]interface{}{
		"title": "Pages",
		"pages": pages,
		"page":  page,
	}

	html, err := h.renderer.Render("pages/index.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Show displays a single page
func (h *PageHandler) Show(ctx router.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.String(http.StatusBadRequest, "Slug is required")
	}

	page, err := h.pageService.GetPageBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Page not found")
	}

	data := map[string]interface{}{
		"title": page.Title,
		"page":  page,
	}

	html, err := h.renderer.Render("pages/show.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// New displays form to create new page
func (h *PageHandler) New(ctx router.Context) error {
	data := map[string]interface{}{
		"title": "New Page",
	}

	html, err := h.renderer.Render("pages/new.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Create handles page creation
func (h *PageHandler) Create(ctx router.Context) error {
	// Get authenticated user
	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return ctx.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Parse form
	if err := ctx.Request().ParseForm(); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid form data")
	}

	title := strings.TrimSpace(ctx.Request().FormValue("title"))
	content := strings.TrimSpace(ctx.Request().FormValue("content"))
	status := ctx.Request().FormValue("status")

	if title == "" || content == "" {
		return ctx.String(http.StatusBadRequest, "Title and content are required")
	}

	// Generate slug
	slug := helpers.GenerateSlug(title)

	// Create page
	page := &domain.Page{
		Title:    title,
		Slug:     slug,
		Content:  content,
		AuthorID: int64(user.ID),
		Status:   domain.PageStatus(status),
	}

	if err := h.pageService.CreatePage(ctx.Request().Context(), page); err != nil {
		log.Printf("Error creating page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error creating page: "+err.Error())
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/pages/%s", page.Slug), http.StatusSeeOther)
	return nil
}

// Edit displays form to edit page
func (h *PageHandler) Edit(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid page ID")
	}

	page, err := h.pageService.GetPageByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Page not found")
	}

	data := map[string]interface{}{
		"title": "Edit Page",
		"page":  page,
	}

	html, err := h.renderer.Render("pages/edit.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Update handles page updates
func (h *PageHandler) Update(ctx router.Context) error {
	// Get authenticated user
	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return ctx.String(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid page ID")
	}

	// Get existing page
	page, err := h.pageService.GetPageByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Page not found")
	}

	// Check authorization
	if page.AuthorID != int64(user.ID) && user.Role != models.RoleAdmin {
		return ctx.String(http.StatusForbidden, "You don't have permission to edit this page")
	}

	// Parse form
	if err := ctx.Request().ParseForm(); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid form data")
	}

	title := strings.TrimSpace(ctx.Request().FormValue("title"))
	content := strings.TrimSpace(ctx.Request().FormValue("content"))
	status := ctx.Request().FormValue("status")

	if title == "" || content == "" {
		return ctx.String(http.StatusBadRequest, "Title and content are required")
	}

	// Update page
	page.Title = title
	page.Content = content
	page.Slug = helpers.GenerateSlug(title)
	if status != "" {
		page.Status = domain.PageStatus(status)
	}

	if err := h.pageService.UpdatePage(ctx.Request().Context(), page); err != nil {
		log.Printf("Error updating page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error updating page")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/pages/%s", page.Slug), http.StatusSeeOther)
	return nil
}

// Delete handles page deletion
func (h *PageHandler) Delete(ctx router.Context) error {
	// Get authenticated user
	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return ctx.String(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid page ID")
	}

	// Get existing page for authorization
	page, err := h.pageService.GetPageByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Page not found")
	}

	// Check authorization
	if page.AuthorID != int64(user.ID) && user.Role != models.RoleAdmin {
		return ctx.String(http.StatusForbidden, "You don't have permission to delete this page")
	}

	if err := h.pageService.DeletePage(ctx.Request().Context(), id); err != nil {
		log.Printf("Error deleting page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error deleting page")
	}

	http.Redirect(ctx.Response(), ctx.Request(), "/pages", http.StatusSeeOther)
	return nil
}

// Publish handles publishing a page
func (h *PageHandler) Publish(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid page ID")
	}

	if err := h.pageService.PublishPage(ctx.Request().Context(), id); err != nil {
		log.Printf("Error publishing page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error publishing page")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/pages/%d/edit", id), http.StatusSeeOther)
	return nil
}

// Unpublish handles unpublishing a page
func (h *PageHandler) Unpublish(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid page ID")
	}

	if err := h.pageService.UnpublishPage(ctx.Request().Context(), id); err != nil {
		log.Printf("Error unpublishing page: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error unpublishing page")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/pages/%d/edit", id), http.StatusSeeOther)
	return nil
}
