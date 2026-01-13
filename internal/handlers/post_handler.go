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

// PostHandler handles post-related requests
type PostHandler struct {
	postService *service.PostService
	renderer    *fith.Engine
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService *service.PostService, renderer *fith.Engine) *PostHandler {
	return &PostHandler{
		postService: postService,
		renderer:    renderer,
	}
}

// Index displays list of posts
func (h *PostHandler) Index(ctx router.Context) error {
	// Parse pagination
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20
	offset := (page - 1) * perPage

	// Get posts
	posts, err := h.postService.ListPublishedPosts(ctx.Request().Context(), perPage, offset)
	if err != nil {
		log.Printf("Error listing posts: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error loading posts")
	}

	data := map[string]interface{}{
		"title": "Posts",
		"posts": posts,
		"page":  page,
	}

	html, err := h.renderer.Render("posts/index.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Show displays a single post
func (h *PostHandler) Show(ctx router.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return ctx.String(http.StatusBadRequest, "Slug is required")
	}

	post, err := h.postService.GetPostBySlug(ctx.Request().Context(), slug)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Post not found")
	}

	data := map[string]interface{}{
		"title": post.Title,
		"post":  post,
	}

	html, err := h.renderer.Render("posts/show.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// New displays form to create new post
func (h *PostHandler) New(ctx router.Context) error {
	data := map[string]interface{}{
		"title": "New Post",
	}

	html, err := h.renderer.Render("posts/new.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Create handles post creation
func (h *PostHandler) Create(ctx router.Context) error {
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

	// Create post
	post := &domain.Post{
		Title:    title,
		Slug:     slug,
		Content:  content,
		AuthorID: int64(user.ID),
		Status:   domain.PostStatus(status),
	}

	if err := h.postService.CreatePost(ctx.Request().Context(), post); err != nil {
		log.Printf("Error creating post: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error creating post: "+err.Error())
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/posts/%s", post.Slug), http.StatusSeeOther)
	return nil
}

// Edit displays form to edit post
func (h *PostHandler) Edit(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid post ID")
	}

	post, err := h.postService.GetPostByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Post not found")
	}

	data := map[string]interface{}{
		"title": "Edit Post",
		"post":  post,
	}

	html, err := h.renderer.Render("posts/edit.html", data)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}

// Update handles post updates
func (h *PostHandler) Update(ctx router.Context) error {
	// Get authenticated user
	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return ctx.String(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid post ID")
	}

	// Get existing post
	post, err := h.postService.GetPostByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Post not found")
	}

	// Check authorization
	if post.AuthorID != int64(user.ID) && user.Role != models.RoleAdmin {
		return ctx.String(http.StatusForbidden, "You don't have permission to edit this post")
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

	// Update post
	post.Title = title
	post.Content = content
	post.Slug = helpers.GenerateSlug(title)
	if status != "" {
		post.Status = domain.PostStatus(status)
	}

	if err := h.postService.UpdatePost(ctx.Request().Context(), post); err != nil {
		log.Printf("Error updating post: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error updating post")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/posts/%s", post.Slug), http.StatusSeeOther)
	return nil
}

// Delete handles post deletion
func (h *PostHandler) Delete(ctx router.Context) error {
	// Get authenticated user
	user, ok := ctx.Get("user").(*models.User)
	if !ok || user == nil {
		return ctx.String(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid post ID")
	}

	// Get existing post for authorization
	post, err := h.postService.GetPostByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.String(http.StatusNotFound, "Post not found")
	}

	// Check authorization
	if post.AuthorID != int64(user.ID) && user.Role != models.RoleAdmin {
		return ctx.String(http.StatusForbidden, "You don't have permission to delete this post")
	}

	if err := h.postService.DeletePost(ctx.Request().Context(), id); err != nil {
		log.Printf("Error deleting post: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error deleting post")
	}

	http.Redirect(ctx.Response(), ctx.Request(), "/posts", http.StatusSeeOther)
	return nil
}

// Publish handles publishing a post
func (h *PostHandler) Publish(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid post ID")
	}

	if err := h.postService.PublishPost(ctx.Request().Context(), id); err != nil {
		log.Printf("Error publishing post: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error publishing post")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/posts/%d/edit", id), http.StatusSeeOther)
	return nil
}

// Unpublish handles unpublishing a post
func (h *PostHandler) Unpublish(ctx router.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid post ID")
	}

	if err := h.postService.UnpublishPost(ctx.Request().Context(), id); err != nil {
		log.Printf("Error unpublishing post: %v", err)
		return ctx.String(http.StatusInternalServerError, "Error unpublishing post")
	}

	http.Redirect(ctx.Response(), ctx.Request(), fmt.Sprintf("/posts/%d/edit", id), http.StatusSeeOther)
	return nil
}
