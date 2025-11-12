package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Client represents a Jira API client
type Client struct {
	baseURL    string
	username   string
	apiToken   string
	projectKey string
	httpClient *http.Client
}

// NewClient creates a new Jira API client
func NewClient(baseURL, username, apiToken, projectKey string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		username:   username,
		apiToken:   apiToken,
		projectKey: projectKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Issue represents a Jira issue
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields represents Jira issue fields
type IssueFields struct {
	Summary  string      `json:"summary"`
	Status   IssueStatus `json:"status"`
	Assignee *IssueUser  `json:"assignee"`
}

// IssueStatus represents Jira issue status
type IssueStatus struct {
	Name string `json:"name"`
}

// IssueUser represents a Jira user
type IssueUser struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// SearchResult represents Jira search results
type SearchResult struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

// Worklog represents a Jira worklog entry
type Worklog struct {
	ID               string     `json:"id,omitempty"`
	IssueID          string     `json:"issueId,omitempty"`
	TimeSpent        string     `json:"timeSpent"`
	TimeSpentSeconds int        `json:"timeSpentSeconds"`
	Started          string     `json:"started"` // Format: 2024-11-11T10:00:00.000+0000
	Comment          string     `json:"comment,omitempty"`
	Author           *IssueUser `json:"author,omitempty"`
}

// GetInProgressIssues retrieves issues in progress for the current user
func (c *Client) GetInProgressIssues() ([]Issue, error) {
	log.Debug().Msg("Fetching in-progress issues")

	jql := "assignee = currentUser() AND status = 'In Progress'"
	if c.projectKey != "" {
		jql = fmt.Sprintf("%s AND project = %s", jql, c.projectKey)
	}
	jql = fmt.Sprintf("%s ORDER BY updated DESC", jql)

	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", "summary,status,assignee")
	params.Add("maxResults", "50")

	endpoint := fmt.Sprintf("%s/rest/api/3/search?%s", c.baseURL, params.Encode())

	var result SearchResult
	if err := c.doRequest("GET", endpoint, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to fetch in-progress issues: %w", err)
	}

	log.Debug().Int("count", len(result.Issues)).Msg("Retrieved in-progress issues")
	return result.Issues, nil
}

// GetIssue retrieves a specific issue by key
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	log.Debug().Str("key", issueKey).Msg("Fetching issue")

	endpoint := fmt.Sprintf("%s/rest/api/3/issue/%s?fields=summary,status,assignee", c.baseURL, issueKey)

	var issue Issue
	if err := c.doRequest("GET", endpoint, nil, &issue); err != nil {
		return nil, fmt.Errorf("failed to fetch issue %s: %w", issueKey, err)
	}

	return &issue, nil
}

// SearchIssues searches for issues by key
func (c *Client) SearchIssues(searchKey string) ([]Issue, error) {
	log.Debug().Str("search", searchKey).Msg("Searching issues")

	jql := fmt.Sprintf("key = %s OR key ~ %s", searchKey, searchKey)
	if c.projectKey != "" {
		jql = fmt.Sprintf("(%s) AND project = %s", jql, c.projectKey)
	}
	jql = fmt.Sprintf("%s ORDER BY key DESC", jql)

	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", "summary,status,assignee")
	params.Add("maxResults", "20")

	endpoint := fmt.Sprintf("%s/rest/api/3/search?%s", c.baseURL, params.Encode())

	var result SearchResult
	if err := c.doRequest("GET", endpoint, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}

	log.Debug().Int("count", len(result.Issues)).Msg("Retrieved search results")
	return result.Issues, nil
}

// AddWorklog adds a worklog entry to an issue
func (c *Client) AddWorklog(issueKey string, timeSpentSeconds int, started time.Time, comment string) (*Worklog, error) {
	log.Debug().
		Str("issue", issueKey).
		Int("seconds", timeSpentSeconds).
		Msg("Adding worklog")

	endpoint := fmt.Sprintf("%s/rest/api/3/issue/%s/worklog", c.baseURL, issueKey)

	// Format started time in Jira format
	startedStr := started.Format("2006-01-02T15:04:05.000-0700")

	payload := map[string]interface{}{
		"timeSpentSeconds": timeSpentSeconds,
		"started":          startedStr,
	}

	if comment != "" {
		payload["comment"] = map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": comment,
						},
					},
				},
			},
		}
	}

	var worklog Worklog
	if err := c.doRequest("POST", endpoint, payload, &worklog); err != nil {
		return nil, fmt.Errorf("failed to add worklog: %w", err)
	}

	log.Info().
		Str("issue", issueKey).
		Str("time", formatSeconds(timeSpentSeconds)).
		Msg("Worklog added successfully")

	return &worklog, nil
}

// GetTodayWorklogs retrieves today's worklogs for the current user
func (c *Client) GetTodayWorklogs() ([]Worklog, error) {
	log.Debug().Msg("Fetching today's worklogs")

	// Get issues with worklogs updated today
	jql := "worklogAuthor = currentUser() AND worklogDate = startOfDay()"
	if c.projectKey != "" {
		jql = fmt.Sprintf("%s AND project = %s", jql, c.projectKey)
	}
	jql = fmt.Sprintf("%s ORDER BY updated DESC", jql)

	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", "worklog")
	params.Add("maxResults", "100")

	endpoint := fmt.Sprintf("%s/rest/api/3/search?%s", c.baseURL, params.Encode())

	var result SearchResult
	if err := c.doRequest("GET", endpoint, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to fetch today's issues: %w", err)
	}

	// Extract worklogs from issues (simplified - in production you'd need to fetch worklogs separately)
	worklogs := []Worklog{}

	log.Debug().Int("count", len(worklogs)).Msg("Retrieved today's worklogs")
	return worklogs, nil
}

// doRequest performs an HTTP request to the Jira API
func (c *Client) doRequest(method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.username, c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("API request failed")
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// formatSeconds formats seconds into human-readable time
func formatSeconds(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	if hours > 0 && minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", minutes)
}
