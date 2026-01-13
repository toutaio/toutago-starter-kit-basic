package middleware

import (
	"net/http"

	"github.com/toutaio/toutago-cosan-router"
)

// RequireOwnership checks if the authenticated user owns the resource
// The resourceOwnerID should be set in the context by the handler
func RequireOwnership(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(c cosan.Context) error {
		user := GetAuthUser(c)
		if user == nil {
			http.Redirect(c.Response(), c.Request(), "/auth/login", http.StatusFound)
			return nil
		}

		// Check if user is admin (admins can access all resources)
		if user.Role == "admin" {
			return next(c)
		}

		// Get resource owner ID from context (should be set by handler)
		resourceOwnerID := c.Get("resource_owner_id")
		if resourceOwnerID == nil {
			c.Response().WriteHeader(http.StatusForbidden)
			c.Response().Write([]byte("Access denied"))
			return nil
		}

		ownerID, ok := resourceOwnerID.(int64)
		if !ok {
			c.Response().WriteHeader(http.StatusInternalServerError)
			c.Response().Write([]byte("Invalid resource owner ID"))
			return nil
		}

		// Check if user owns the resource
		if int64(user.ID) != ownerID {
			c.Response().WriteHeader(http.StatusForbidden)
			c.Response().Write([]byte("You don't have permission to access this resource"))
			return nil
		}

		return next(c)
	}
}

// CanEdit checks if user can edit a resource
// Editors and admins can edit any content, users can only edit their own
func CanEdit(userID int64, userRole string, resourceOwnerID int64) bool {
	if userRole == "admin" || userRole == "editor" {
		return true
	}
	return userID == resourceOwnerID
}

// CanDelete checks if user can delete a resource
// Only admins and resource owners can delete
func CanDelete(userID int64, userRole string, resourceOwnerID int64) bool {
	if userRole == "admin" {
		return true
	}
	return userID == resourceOwnerID
}

// CanPublish checks if user can publish content
// Only admins and editors can publish
func CanPublish(userRole string) bool {
	return userRole == "admin" || userRole == "editor"
}

// CanManageUsers checks if user can manage other users
// Only admins can manage users
func CanManageUsers(userRole string) bool {
	return userRole == "admin"
}
