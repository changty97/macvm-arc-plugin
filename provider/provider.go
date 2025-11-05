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

// MacVMProvider manages communication with macvmagt
type MacVMProvider struct {
	AgentURL string
	Client   *http.Client
}

// VMProvisionRequest mirrors your Bash payload
type VMProvisionRequest struct {
	VMID                    string   `json:"vmId"`
	ImageName               string   `json:"imageName"`
	RunnerRegistrationToken string   `json:"runnerRegistrationToken"`
	RunnerName              string   `json:"runnerName"`
	RunnerLabels            []string `json:"runnerLabels"`
}

// NewMacVMProvider initializes the provider
func NewMacVMProvider(agentURL string) *MacVMProvider {
	return &MacVMProvider{
		AgentURL: agentURL,
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CreateRunner provisions a new macOS VM runner
func (p *MacVMProvider) CreateRunner(ctx context.Context, runnerName string, token string) error {
	vmID := fmt.Sprintf("%s-%d", runnerName, time.Now().Unix())

	payload := VMProvisionRequest{
		VMID:                    vmID,
		ImageName:               "macos_26",
		RunnerRegistrationToken: token,
		RunnerName:              fmt.Sprintf("github-runner-%s", vmID),
		RunnerLabels:            []string{"macos", "arm64"},
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/provision-vm", p.AgentURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call agent: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("[macvm-arc-plugin] Provisioned VM %s via %s\n", vmID, p.AgentURL)
	return nil
}

// DeleteRunner tears down a VM runner
func (p *MacVMProvider) DeleteRunner(ctx context.Context, vmID string) error {
	payload := map[string]string{"vmId": vmID}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/delete-vm", p.AgentURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to build delete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call delete endpoint: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected delete status %d: %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("[macvm-arc-plugin] Deleted VM %s\n", vmID)
	return nil
}

// ListVMs (optional helper)
func (p *MacVMProvider) ListVMs(ctx context.Context) error {
	url := fmt.Sprintf("%s/vms", p.AgentURL)
	resp, err := p.Client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to list vms: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("[macvm-arc-plugin] Current VM state:\n", string(body))
	return nil
}
