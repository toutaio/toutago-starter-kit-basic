// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	router "github.com/toutaio/toutago-cosan-router"
)

// Logger logs HTTP requests with timing information.
func Logger(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.Context) error {
		start := time.Now()

		// Call next handler
		err := next(ctx)

		// Log request details
		duration := time.Since(start)
		log.Printf(
			"%s %s - Duration: %v",
			ctx.Request().Method,
			ctx.Request().URL.Path,
			duration,
		)

		return err
	}
}

// Recovery recovers from panics and returns a 500 error.
func Recovery(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC: %v", r)
				// Try to send error response
				if err := ctx.String(http.StatusInternalServerError, "Internal Server Error"); err != nil {
					log.Printf("Failed to send error response: %v", err)
				}
			}
		}()

		return next(ctx)
	}
}

// SecurityHeaders adds security-related HTTP headers to responses.
func SecurityHeaders(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.Context) error {
		// Set security headers
		ctx.Response().Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Response().Header().Set("X-Frame-Options", "DENY")
		ctx.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		ctx.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return next(ctx)
	}
}

// CORS adds Cross-Origin Resource Sharing headers.
func CORS(allowedOrigins []string) func(router.HandlerFunc) router.HandlerFunc {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(ctx router.Context) error {
			origin := ctx.Request().Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, ao := range allowedOrigins {
				if ao == "*" || ao == origin {
					allowed = true
					break
				}
			}

			if allowed {
				ctx.Response().Header().Set("Access-Control-Allow-Origin", origin)
				ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				ctx.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight request
			if ctx.Request().Method == http.MethodOptions {
				return ctx.String(http.StatusOK, "")
			}

			return next(ctx)
		}
	}
}

// RequestID adds a unique request ID to each request.
func RequestID(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx router.Context) error {
		requestID := ctx.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate simple request ID
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		ctx.Response().Header().Set("X-Request-ID", requestID)
		ctx.Set("request_id", requestID)

		return next(ctx)
	}
}
