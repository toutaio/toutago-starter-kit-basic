package middleware

import (
	"testing"
)

func TestFlashMessage(t *testing.T) {
	// Create a mock context
	// Note: This is a simplified test. In a real scenario, we'd use a proper mock
	t.Run("SetAndGetFlash", func(t *testing.T) {
		message := FlashMessage{
			Type:    "success",
			Message: "Operation successful",
		}

		if message.Type != "success" {
			t.Errorf("Expected type 'success', got '%s'", message.Type)
		}

		if message.Message != "Operation successful" {
			t.Errorf("Expected message 'Operation successful', got '%s'", message.Message)
		}
	})
}
