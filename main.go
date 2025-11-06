package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/changty97/macvm-arc-plugin/provider"
)

func main() {
	agentURL := os.Getenv("AGENT_URL")
	if agentURL == "" {
		agentURL = "http://localhost:8081"
		fmt.Printf("AGENT_URL not set, falling back to: %s", agentURL)
	}
	fmt.Printf("AGENT_URL set to: %s", agentURL)

	p := provider.NewMacVMProvider(agentURL)

	// Use a reasonable context timeout for the entire sequence
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	runnerID := "test-runner"
	token := "fake-token-runner"

	fmt.Println("--- Starting Agent Client Workflow ---")

	// 1. List Current VMs
	fmt.Println("1. Checking current VM status...")
	p.ListVMs(ctx)

	// 2. Create Runner
	fmt.Printf("2. Sending request to create runner '%s'...", runnerID)
	// The agent responds quickly (202 Accepted) because provisioning is asynchronous.
	err := p.CreateRunner(ctx, runnerID, token)
	if err != nil {
		fmt.Errorf("Failed to create runner: %v", err)
	}
	fmt.Printf("Successfully requested runner creation for '%s'.", runnerID)

	// 3. Wait/Monitor (Simulated Wait)
	// NOTE: In a real ARC orchestrator, this static time.Sleep would be replaced
	// by a dedicated monitoring loop that checks the agent's /vms endpoint
	// until the VM is reported as 'ready' (e.g., connected to GitHub ARC).
	fmt.Printf("3. Waiting for 5 seconds to simulate GitHub runner connection...")
	time.Sleep(5 * time.Second)
	fmt.Println("... Wait complete.")

	// Check status after simulated run
	fmt.Println("3b. Checking status after provisioning wait...")
	p.ListVMs(ctx)

	// 4. Delete Runner
	fmt.Printf("4. Sending request to delete runner '%s'...", runnerID)
	// The agent handles process termination (kill -9) and directory removal asynchronously.
	err = p.DeleteRunner(ctx, runnerID)
	if err != nil {
		fmt.Errorf("Failed to delete runner: %v", err)
	}
	fmt.Printf("Successfully requested runner deletion for '%s'.", runnerID)

	// 5. Final Check
	fmt.Println("5. Checking final VM status (wait for deletion goroutine)...")
	// Give the deletion goroutine time to work before the client exits
	time.Sleep(2 * time.Second)
	p.ListVMs(ctx)

	fmt.Println("--- Agent Client Workflow Complete ---")
}
