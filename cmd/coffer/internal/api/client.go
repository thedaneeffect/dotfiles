package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for the Cloudflare Worker backend
type Client struct {
	BaseURL    string
	Passphrase string
	HTTPClient *http.Client
}

// Metadata contains information about uploaded files
type Metadata struct {
	Files []string `json:"files"`
	Size  string   `json:"size"`
}

// GroupMetadata contains metadata for a group stored in the worker
type GroupMetadata struct {
	Files    []string  `json:"files"`
	Size     string    `json:"size"`
	Uploaded time.Time `json:"uploaded"`
}

// NewClient creates a new API client
func NewClient(baseURL, passphrase string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Passphrase: passphrase,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Push uploads encrypted files to the worker
func (c *Client) Push(group string, data string, metadata Metadata) error {
	url := fmt.Sprintf("%s/secrets/%s", c.BaseURL, group)

	// Create request with base64 data as body
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.Passphrase)
	req.Header.Set("Content-Type", "text/plain")

	// Add metadata headers
	filesJSON, err := json.Marshal(metadata.Files)
	if err != nil {
		return fmt.Errorf("failed to marshal files: %w", err)
	}
	req.Header.Set("X-Files", string(filesJSON))
	req.Header.Set("X-Size", metadata.Size)

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized - check COFFER_PASSPHRASE")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden - invalid credentials")
	case http.StatusOK:
		// Success
	default:
		return fmt.Errorf("push failed (HTTP %d)", resp.StatusCode)
	}

	return nil
}

// Pull downloads encrypted files from the worker
func (c *Client) Pull(group string) (string, error) {
	url := fmt.Sprintf("%s/secrets/%s", c.BaseURL, group)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Passphrase)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return "", fmt.Errorf("no files found for group: %s", group)
	case http.StatusUnauthorized:
		return "", fmt.Errorf("unauthorized - check COFFER_PASSPHRASE")
	case http.StatusForbidden:
		return "", fmt.Errorf("forbidden - invalid credentials")
	case http.StatusOK:
		// Success, continue below
	default:
		return "", fmt.Errorf("pull failed (HTTP %d)", resp.StatusCode)
	}

	// Read base64 data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(data), nil
}

// Delete removes a group from the worker
func (c *Client) Delete(group string) error {
	url := fmt.Sprintf("%s/secrets/%s", c.BaseURL, group)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Passphrase)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("group not found in worker: %s", group)
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized - check COFFER_PASSPHRASE")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden - invalid credentials")
	case http.StatusOK:
		// Success
	default:
		return fmt.Errorf("delete failed (HTTP %d)", resp.StatusCode)
	}

	return nil
}

// GetMetadata fetches metadata for all groups from the worker
func (c *Client) GetMetadata() (map[string]GroupMetadata, error) {
	url := fmt.Sprintf("%s/metadata", c.BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Passphrase)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("unauthorized - check COFFER_PASSPHRASE")
	case http.StatusForbidden:
		return nil, fmt.Errorf("forbidden - invalid credentials")
	case http.StatusNotFound:
		return nil, fmt.Errorf("metadata endpoint not found")
	case http.StatusOK:
		// Success, continue below
	default:
		return nil, fmt.Errorf("get metadata failed (HTTP %d)", resp.StatusCode)
	}

	var metadata map[string]GroupMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return metadata, nil
}

// ListGroups returns a list of all group names stored in the worker
func (c *Client) ListGroups() ([]string, error) {
	metadata, err := c.GetMetadata()
	if err != nil {
		return nil, err
	}

	groups := make([]string, 0, len(metadata))
	for group := range metadata {
		groups = append(groups, group)
	}

	return groups, nil
}
