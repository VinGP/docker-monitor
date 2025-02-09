package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"pinger/internal/model"
	"time"
)

const requestTimeout = 5 * time.Second

type BackendAPI struct {
	Client  *http.Client
	BaseURL string
}

// NewBackendAPI initializes the BackendAPI with an optional custom HTTP client.
func NewBackendAPI(baseURL string) *BackendAPI {
	return &BackendAPI{
		BaseURL: baseURL,
		Client:  &http.Client{}, // Using default HTTP client
	}
}

// SaveContainerStatus sends the container status to the backend and returns the response or error message.
func (api *BackendAPI) SaveContainerStatus(containerStatus model.ContainerStatus) error {
	url := api.BaseURL + "/container/status"

	payload, err := json.Marshal(containerStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal container status: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := api.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Error closing body", slog.Any("error", err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned non-OK status: %s", resp.Status) //nolint:err113 //nolint
	}

	var response struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("backend error: %s", response.Error) //nolint:err113 //nolint
	}

	return nil
}
