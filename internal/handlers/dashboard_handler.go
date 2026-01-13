package handlers

import (
	"log"
	"net/http"
	"time"

	router "github.com/toutaio/toutago-cosan-router"
	"github.com/toutaio/toutago-fith-renderer"
	"github.com/toutaio/toutago-starter-kit-basic/internal/domain"
	"github.com/toutaio/toutago-starter-kit-basic/internal/models"
	"github.com/toutaio/toutago-starter-kit-basic/internal/repository"
)

type DashboardHandler struct {
	renderer *fith.Engine
	postRepo *repository.PostRepository
	pageRepo *repository.PageRepository
}

func NewDashboardHandler(renderer *fith.Engine, postRepo *repository.PostRepository, pageRepo *repository.PageRepository) *DashboardHandler {
	return &DashboardHandler{
		renderer: renderer,
		postRepo: postRepo,
		pageRepo: pageRepo,
	}
}

type ActivityItem struct {
	Type      string
	Title     string
	Action    string
	Timestamp time.Time
}

func (h *DashboardHandler) Show(c router.Context) error {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusSeeOther)
		return nil
	}

	// Get user's content statistics
	var postCount, pageCount int
	var recentPosts []*domain.Post
	var recentPages []*domain.Page

	if h.postRepo != nil {
		posts, err := h.postRepo.ListByAuthor(c.Request().Context(), int64(user.ID), 5, 0)
		if err == nil {
			postCount = len(posts)
			recentPosts = posts
		} else {
			log.Printf("Error fetching posts: %v", err)
		}
	}

	if h.pageRepo != nil {
		pages, err := h.pageRepo.ListByAuthor(c.Request().Context(), int64(user.ID), 5, 0)
		if err == nil {
			pageCount = len(pages)
			recentPages = pages
		} else {
			log.Printf("Error fetching pages: %v", err)
		}
	}

	// Build recent activity feed
	activity := buildActivityFeed(recentPosts, recentPages)

	data := map[string]interface{}{
		"User":           user,
		"PostCount":      postCount,
		"PageCount":      pageCount,
		"RecentPosts":    recentPosts,
		"RecentPages":    recentPages,
		"RecentActivity": activity,
	}

	html, err := h.renderer.Render("pages/dashboard.html", data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error rendering template: "+err.Error())
	}

	return c.HTML(http.StatusOK, html)
}

func buildActivityFeed(posts []*domain.Post, pages []*domain.Page) []ActivityItem {
	var activity []ActivityItem

	for _, post := range posts {
		activity = append(activity, ActivityItem{
			Type:      "post",
			Title:     post.Title,
			Action:    "created",
			Timestamp: post.CreatedAt,
		})
	}

	for _, page := range pages {
		activity = append(activity, ActivityItem{
			Type:      "page",
			Title:     page.Title,
			Action:    "created",
			Timestamp: page.CreatedAt,
		})
	}

	// Sort by timestamp (most recent first)
	// Simple bubble sort for small lists
	for i := 0; i < len(activity); i++ {
		for j := i + 1; j < len(activity); j++ {
			if activity[j].Timestamp.After(activity[i].Timestamp) {
				activity[i], activity[j] = activity[j], activity[i]
			}
		}
	}

	// Limit to 10 items
	if len(activity) > 10 {
		activity = activity[:10]
	}

	return activity
}
