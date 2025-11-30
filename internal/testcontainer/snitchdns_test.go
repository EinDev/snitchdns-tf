package testcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestSnitchDNSContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start container
	container, err := NewSnitchDNSContainer(ctx, SnitchDNSContainerRequest{
		ExposePorts: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	t.Logf("Container started successfully")
	t.Logf("HTTP Host: %s", container.HTTPHost)
	t.Logf("API Key: %s", container.APIKey[:10]+"...")

	// Test API connectivity
	t.Run("API Authentication", func(t *testing.T) {
		testAPIAuthentication(t, container)
	})

	t.Run("API Zones Endpoint", func(t *testing.T) {
		testZonesEndpoint(t, container)
	})

	t.Run("API Record Types", func(t *testing.T) {
		testRecordTypes(t, container)
	})
}

func testAPIAuthentication(t *testing.T, container *SnitchDNSContainer) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Test without API key (should fail)
	req, err := http.NewRequest("GET", container.GetAPIEndpoint()+"/zones", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 without API key, got %d", resp.StatusCode)
	}

	// Test with API key (should succeed)
	req, err = http.NewRequest("GET", container.GetAPIEndpoint()+"/zones", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-SnitchDNS-Auth", container.APIKey)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200 with API key, got %d: %s", resp.StatusCode, string(body))
	}

	t.Logf("✓ API authentication working correctly")
}

func testZonesEndpoint(t *testing.T, container *SnitchDNSContainer) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", container.GetAPIEndpoint()+"/zones", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-SnitchDNS-Auth", container.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check for pagination fields
	if _, ok := result["page"]; !ok {
		t.Error("Response missing 'page' field")
	}
	if _, ok := result["data"]; !ok {
		t.Error("Response missing 'data' field")
	}

	t.Logf("✓ Zones endpoint returning valid response")
	t.Logf("  Response: %+v", result)
}

func testRecordTypes(t *testing.T, container *SnitchDNSContainer) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", container.GetAPIEndpoint()+"/records/types", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-SnitchDNS-Auth", container.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var types []string
	if err := json.Unmarshal(body, &types); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	expectedTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT"}
	for _, expectedType := range expectedTypes {
		found := false
		for _, t := range types {
			if t == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected record type '%s' not found in response", expectedType)
		}
	}

	t.Logf("✓ Record types endpoint working correctly")
	t.Logf("  Found %d record types: %v", len(types), types)
}

// BenchmarkContainerStartup benchmarks container startup time
func BenchmarkContainerStartup(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		container, err := NewSnitchDNSContainer(ctx, SnitchDNSContainerRequest{
			ExposePorts: true,
		})
		if err != nil {
			b.Fatalf("Failed to start container: %v", err)
		}

		if err := container.Terminate(ctx); err != nil {
			b.Logf("Failed to terminate container: %v", err)
		}
	}
}

// Example usage for documentation
func ExampleNewSnitchDNSContainer() {
	ctx := context.Background()

	container, err := NewSnitchDNSContainer(ctx, SnitchDNSContainerRequest{
		ExposePorts: true,
	})
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	// Use the container
	fmt.Printf("API Endpoint: %s\n", container.GetAPIEndpoint())
	fmt.Printf("API Key: %s\n", container.APIKey)
}
