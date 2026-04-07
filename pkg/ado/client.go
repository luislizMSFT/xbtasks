package ado

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an Azure DevOps REST API client for a single org/project.
type Client struct {
	org        string
	project    string
	token      string
	httpClient *http.Client
}

// NewClient creates an ADO client for the given organization, project, and Bearer token.
func NewClient(org, project, token string) *Client {
	return &Client{
		org:     org,
		project: project,
		token:   token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// NewMultiOrgClients creates one Client per org/project pair, all using the same token.
func NewMultiOrgClients(orgProjects []OrgProject, token string) []*Client {
	clients := make([]*Client, len(orgProjects))
	for i, op := range orgProjects {
		clients[i] = NewClient(op.Org, op.Project, token)
	}
	return clients
}

// Org returns the organization name for this client.
func (c *Client) Org() string { return c.org }

// Project returns the project name for this client.
func (c *Client) Project() string { return c.project }

// apiURL constructs the full ADO REST API URL for a given path.
// Example: apiURL("wit/wiql") → "https://dev.azure.com/{org}/{project}/_apis/wit/wiql"
func (c *Client) apiURL(path string) string {
	return fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/%s", c.org, c.project, path)
}

// doRequest executes an HTTP request with Bearer auth and the specified content type.
func (c *Client) doRequest(method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}

	return resp, nil
}

// doJSON marshals body to JSON and executes a request with application/json content type.
func (c *Client) doJSON(method, url string, body any) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON: %w", err)
	}
	return c.doRequest(method, url, bytes.NewReader(jsonBytes), "application/json")
}

// doPatch marshals patch operations and executes with application/json-patch+json content type.
// CRITICAL: ADO requires this specific content type for JSON Patch operations.
func (c *Client) doPatch(url string, ops []PatchOperation) (*http.Response, error) {
	jsonBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, fmt.Errorf("marshal patch ops: %w", err)
	}
	return c.doRequest("PATCH", url, bytes.NewReader(jsonBytes), "application/json-patch+json")
}

// decodeResponse reads the response body and unmarshals JSON into target.
// Returns an error if the status code indicates failure.
func decodeResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ADO API returned %d: %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
