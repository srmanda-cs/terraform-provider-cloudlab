// Package provider implements the CloudLab Terraform provider.
package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultPortalURL = "https://www.cloudlab.us"
	apiBasePath      = "/portal/api"

	// Poll intervals for experiment readiness.
	pollInterval = 15 * time.Second
	pollTimeout  = 30 * time.Minute
)

// Experiment status strings returned by the Portal API.
const (
	StatusReady    = "ready"
	StatusFailed   = "failed"
	StatusCreating = "created"
)

// Client is the CloudLab Portal API HTTP client.
type Client struct {
	httpClient *http.Client
	portalURL  string
	token      string
}

// NewClient creates a new CloudLab Portal API client.
func NewClient(portalURL, token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		portalURL:  portalURL,
		token:      token,
	}
}

// ExperimentCreateRequest represents the body for POST /experiments.
type ExperimentCreateRequest struct {
	Name           string         `json:"name"`
	Project        string         `json:"project"`
	ProfileName    string         `json:"profile_name"`
	ProfileProject string         `json:"profile_project"`
	Duration       *int64         `json:"duration,omitempty"`
	StartAt        *string        `json:"start_at,omitempty"`
	StopAt         *string        `json:"stop_at,omitempty"`
	Bindings       map[string]any `json:"bindings,omitempty"`
}

// ExperimentResponse represents a CloudLab experiment object.
type ExperimentResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Project        string  `json:"project"`
	Group          string  `json:"group"`
	ProfileID      string  `json:"profile_id"`
	ProfileName    string  `json:"profile_name"`
	ProfileProject string  `json:"profile_project"`
	Creator        string  `json:"creator"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
	StartAt        *string `json:"start_at"`
	StopAt         *string `json:"stop_at"`
	ExpiresAt      *string `json:"expires_at"`
}

// ExperimentListResponse represents the GET /experiments response.
type ExperimentListResponse struct {
	Experiments []ExperimentResponse `json:"experiments"`
}

// ProfileCreateRequest represents the body for POST /profiles.
type ProfileCreateRequest struct {
	Name            string `json:"name"`
	Project         string `json:"project"`
	Script          string `json:"script,omitempty"`
	RepositoryURL   string `json:"repository_url,omitempty"`
	Public          bool   `json:"public"`
	ProjectWritable bool   `json:"project_writable"`
}

// ProfileResponse represents a CloudLab profile object.
type ProfileResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Version         int64   `json:"version"`
	Project         string  `json:"project"`
	Creator         string  `json:"creator"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       *string `json:"updated_at"`
	RepositoryURL   *string `json:"repository_url"`
	Public          bool    `json:"public"`
	ProjectWritable bool    `json:"project_writable"`
}

// ProfileListResponse represents the GET /profiles response.
type ProfileListResponse struct {
	Profiles []ProfileResponse `json:"profiles"`
}

// ResgroupCreateRequest represents the body for POST /resgroups.
type ResgroupCreateRequest struct {
	Project   string         `json:"project"`
	Reason    string         `json:"reason"`
	StartAt   *string        `json:"start_at,omitempty"`
	ExpiresAt *string        `json:"expires_at,omitempty"`
	NodeTypes []ResgroupNode `json:"nodetypes,omitempty"`
}

// ResgroupNode represents a node type entry in a reservation group.
type ResgroupNode struct {
	NodeType  string `json:"type"`
	Aggregate string `json:"aggregate"`
	Count     int64  `json:"count"`
}

// ResgroupResponse represents a CloudLab reservation group object.
type ResgroupResponse struct {
	ID        string  `json:"id"`
	Project   string  `json:"project"`
	Reason    string  `json:"reason"`
	Creator   string  `json:"creator"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	StartAt   *string `json:"start_at"`
	ExpiresAt *string `json:"expires_at"`
}

// ResgroupListResponse represents the GET /resgroups response.
type ResgroupListResponse struct {
	Resgroups []ResgroupResponse `json:"resgroups"`
}

// ManifestListResponse represents the GET /experiments/{id}/manifests response.
type ManifestListResponse struct {
	Manifests []ManifestEntry `json:"manifests"`
}

// ManifestEntry represents a single manifest for an aggregate.
type ManifestEntry struct {
	Aggregate string      `json:"aggregate"`
	Nodes     []NodeEntry `json:"nodes"`
}

// NodeEntry represents a node in a manifest.
type NodeEntry struct {
	ClientID   string          `json:"client_id"`
	Hostname   string          `json:"hostname"`
	Interfaces []NodeInterface `json:"interfaces"`
}

// NodeInterface represents a network interface on a node.
type NodeInterface struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// APIError represents an error response from the Portal API.
type APIError struct {
	StatusCode int
	Message    string   `json:"error"`
	Errors     []string `json:"errors"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("cloudlab api error (HTTP %d): %s: %v", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("cloudlab api error (HTTP %d): %s", e.StatusCode, e.Message)
}

func (c *Client) doRequest(method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	url := c.portalURL + apiBasePath + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if jsonErr := json.Unmarshal(respBody, apiErr); jsonErr != nil {
			apiErr.Message = string(respBody)
		}
		return nil, apiErr
	}

	return respBody, nil
}

// CreateExperiment creates a new experiment on CloudLab.
func (c *Client) CreateExperiment(req *ExperimentCreateRequest) (*ExperimentResponse, error) {
	body, err := c.doRequest(http.MethodPost, "/experiments", req)
	if err != nil {
		return nil, err
	}
	var exp ExperimentResponse
	if err := json.Unmarshal(body, &exp); err != nil {
		return nil, fmt.Errorf("failed to parse experiment response: %w", err)
	}
	return &exp, nil
}

// GetExperiment retrieves an experiment by its Portal ID (name:project or UUID).
func (c *Client) GetExperiment(experimentID string) (*ExperimentResponse, error) {
	body, err := c.doRequest(http.MethodGet, "/experiments/"+experimentID, nil)
	if err != nil {
		return nil, err
	}
	var exp ExperimentResponse
	if err := json.Unmarshal(body, &exp); err != nil {
		return nil, fmt.Errorf("failed to parse experiment response: %w", err)
	}
	return &exp, nil
}

// ListExperiments lists experiments, optionally filtered by project.
func (c *Client) ListExperiments(project string) ([]ExperimentResponse, error) {
	path := "/experiments"
	if project != "" {
		path += "?project=" + project
	}
	body, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var list ExperimentListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("failed to parse experiment list response: %w", err)
	}
	return list.Experiments, nil
}

// DeleteExperiment deletes (terminates) an experiment.
func (c *Client) DeleteExperiment(experimentID string) error {
	_, err := c.doRequest(http.MethodDelete, "/experiments/"+experimentID, nil)
	return err
}

// WaitForExperiment polls until the experiment reaches "ready" or "failed" status.
func (c *Client) WaitForExperiment(ctx context.Context, experimentID string) (*ExperimentResponse, error) {
	deadline := time.Now().Add(pollTimeout)
	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for experiment %s to become ready", experimentID)
		}

		exp, err := c.GetExperiment(experimentID)
		if err != nil {
			return nil, fmt.Errorf("error polling experiment status: %w", err)
		}

		switch exp.Status {
		case StatusReady:
			return exp, nil
		case StatusFailed:
			return nil, fmt.Errorf("experiment %s failed to start", experimentID)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// GetManifests retrieves the manifests for a running experiment.
func (c *Client) GetManifests(experimentID string) ([]ManifestEntry, error) {
	body, err := c.doRequest(http.MethodGet, "/experiments/"+experimentID+"/manifests", nil)
	if err != nil {
		return nil, err
	}
	var list ManifestListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("failed to parse manifest response: %w", err)
	}
	return list.Manifests, nil
}

// CreateProfile creates a new experiment profile.
func (c *Client) CreateProfile(req *ProfileCreateRequest) (*ProfileResponse, error) {
	body, err := c.doRequest(http.MethodPost, "/profiles", req)
	if err != nil {
		return nil, err
	}
	var profile ProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile response: %w", err)
	}
	return &profile, nil
}

// GetProfile retrieves a profile by its Portal ID.
func (c *Client) GetProfile(profileID string) (*ProfileResponse, error) {
	body, err := c.doRequest(http.MethodGet, "/profiles/"+profileID, nil)
	if err != nil {
		return nil, err
	}
	var profile ProfileResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse profile response: %w", err)
	}
	return &profile, nil
}

// ListProfiles lists all profiles, optionally filtered by project.
func (c *Client) ListProfiles(project string) ([]ProfileResponse, error) {
	path := "/profiles"
	if project != "" {
		path += "?project=" + project
	}
	body, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var list ProfileListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("failed to parse profile list response: %w", err)
	}
	return list.Profiles, nil
}

// DeleteProfile deletes a profile.
func (c *Client) DeleteProfile(profileID string) error {
	_, err := c.doRequest(http.MethodDelete, "/profiles/"+profileID, nil)
	return err
}

// CreateResgroup creates a new reservation group.
func (c *Client) CreateResgroup(req *ResgroupCreateRequest) (*ResgroupResponse, error) {
	body, err := c.doRequest(http.MethodPost, "/resgroups", req)
	if err != nil {
		return nil, err
	}
	var rg ResgroupResponse
	if err := json.Unmarshal(body, &rg); err != nil {
		return nil, fmt.Errorf("failed to parse resgroup response: %w", err)
	}
	return &rg, nil
}

// GetResgroup retrieves a reservation group by ID.
func (c *Client) GetResgroup(resgroupID string) (*ResgroupResponse, error) {
	body, err := c.doRequest(http.MethodGet, "/resgroups/"+resgroupID, nil)
	if err != nil {
		return nil, err
	}
	var rg ResgroupResponse
	if err := json.Unmarshal(body, &rg); err != nil {
		return nil, fmt.Errorf("failed to parse resgroup response: %w", err)
	}
	return &rg, nil
}

// DeleteResgroup deletes a reservation group.
func (c *Client) DeleteResgroup(resgroupID string) error {
	_, err := c.doRequest(http.MethodDelete, "/resgroups/"+resgroupID, nil)
	return err
}
