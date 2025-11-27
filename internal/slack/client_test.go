package slack

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	userToken := "xoxp-test-token"
	channelID := "C1234567890"

	client := NewClient(userToken, channelID)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.userToken != userToken {
		t.Errorf("Expected userToken %s, got %s", userToken, client.userToken)
	}

	if client.channelID != channelID {
		t.Errorf("Expected channelID %s, got %s", channelID, client.channelID)
	}

	if client.httpClient == nil {
		t.Error("Expected non-nil HTTP client")
	}
}

func TestClient_Structure(t *testing.T) {
	client := NewClient("token", "channel")

	// Verify client has required methods
	t.Run("SetStatus method exists", func(t *testing.T) {
		err := client.SetStatus("test", ":coffee:", 10)
		// Will fail due to invalid token, but method should exist
		if err == nil {
			t.Skip("Skipping API call test - requires valid credentials")
		}
	})

	t.Run("PostMessage method exists", func(t *testing.T) {
		err := client.PostMessage("test message")
		// Will fail due to invalid token, but method should exist
		if err == nil {
			t.Skip("Skipping API call test - requires valid credentials")
		}
	})

	t.Run("ClearStatus method exists", func(t *testing.T) {
		err := client.ClearStatus()
		// Will fail due to invalid token, but method should exist
		if err == nil {
			t.Skip("Skipping API call test - requires valid credentials")
		}
	})
}
