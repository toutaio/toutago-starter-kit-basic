package middleware

import (
	"github.com/toutaio/toutago-cosan-router"
)

const (
	FlashContextKey = "flash_messages"
	SessionFlashKey = "flash"
)

type FlashMessage struct {
	Type    string // success, error, warning, info
	Message string
}

// FlashMiddleware handles flash messages in the session
func FlashMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
	return func(c cosan.Context) error {
		// Get flash messages from session
		sessionData := c.Get("session")
		if sessionData != nil {
			if session, ok := sessionData.(map[string]interface{}); ok {
				if flashData, exists := session[SessionFlashKey]; exists {
					if messages, ok := flashData.([]FlashMessage); ok {
						c.Set(FlashContextKey, messages)
						// Clear flash messages from session after retrieving
						delete(session, SessionFlashKey)
					}
				}
			}
		}

		return next(c)
	}
}

// SetFlash adds a flash message to the session
func SetFlash(c cosan.Context, flashType, message string) {
	sessionData := c.Get("session")
	if sessionData == nil {
		sessionData = make(map[string]interface{})
		c.Set("session", sessionData)
	}

	if session, ok := sessionData.(map[string]interface{}); ok {
		var messages []FlashMessage
		if existingFlash, exists := session[SessionFlashKey]; exists {
			if existing, ok := existingFlash.([]FlashMessage); ok {
				messages = existing
			}
		}
		messages = append(messages, FlashMessage{
			Type:    flashType,
			Message: message,
		})
		session[SessionFlashKey] = messages
	}
}

// GetFlashMessages retrieves flash messages from context
func GetFlashMessages(c cosan.Context) []FlashMessage {
	if messages := c.Get(FlashContextKey); messages != nil {
		if flash, ok := messages.([]FlashMessage); ok {
			return flash
		}
	}
	return nil
}
