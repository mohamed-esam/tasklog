package jira

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://example.atlassian.net", "user@example.com", "token123", "PROJ")

	if client == nil {
		t.Fatal("expected client to be created")
	}

	if client.baseURL != "https://example.atlassian.net" {
		t.Errorf("expected baseURL to be set correctly, got %s", client.baseURL)
	}

	if client.username != "user@example.com" {
		t.Errorf("expected username to be set correctly")
	}

	if client.apiToken != "token123" {
		t.Errorf("expected apiToken to be set correctly")
	}

	if client.projectKey != "PROJ" {
		t.Errorf("expected projectKey to be set correctly, got %s", client.projectKey)
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected timeout to be 30s, got %v", client.httpClient.Timeout)
	}
}

func TestNewClient_TrimTrailingSlash(t *testing.T) {
	client := NewClient("https://example.atlassian.net/", "user@example.com", "token123", "PROJ")

	if client.baseURL != "https://example.atlassian.net" {
		t.Errorf("expected trailing slash to be trimmed, got %s", client.baseURL)
	}
}

func TestFormatSeconds(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{3600, "1h"},
		{1800, "30m"},
		{5400, "1h 30m"},
		{7200, "2h"},
		{300, "5m"},
		{0, "0m"},
		{7260, "2h 1m"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatSeconds(tt.seconds)
			if result != tt.expected {
				t.Errorf("formatSeconds(%d) = %s, want %s", tt.seconds, result, tt.expected)
			}
		})
	}
}

func TestIssueStructures(t *testing.T) {
	// Test that the structures can be created and marshaled
	issue := Issue{
		Key: "PROJ-123",
		Fields: IssueFields{
			Summary: "Test issue",
			Status: IssueStatus{
				Name: "In Progress",
			},
			Assignee: &IssueUser{
				DisplayName:  "John Doe",
				EmailAddress: "john@example.com",
			},
		},
	}

	if issue.Key != "PROJ-123" {
		t.Error("issue key not set correctly")
	}

	if issue.Fields.Summary != "Test issue" {
		t.Error("issue summary not set correctly")
	}

	if issue.Fields.Status.Name != "In Progress" {
		t.Error("issue status not set correctly")
	}

	if issue.Fields.Assignee == nil {
		t.Error("assignee should not be nil")
	}
}

func TestWorklogStructure(t *testing.T) {
	worklog := Worklog{
		ID:               "12345",
		TimeSpent:        "2h",
		TimeSpentSeconds: 7200,
		Started:          "2024-11-11T10:00:00.000+0000",
		Comment:          []byte(`"Test comment"`),
	}

	if worklog.ID != "12345" {
		t.Error("worklog ID not set correctly")
	}

	if worklog.TimeSpentSeconds != 7200 {
		t.Error("time spent seconds not set correctly")
	}
}
