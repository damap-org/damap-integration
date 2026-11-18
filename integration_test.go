package damap_integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/damap-org/damap-integration"
	"github.com/damap-org/damap-integration/domain"
)

func TestIntegration(t *testing.T) {
	t.Logf("⚙️ Creating integration server...")
	srv, err := damap_integration.New(
		damap_integration.WithListen("127.0.0.1:0"),
	)
	if err != nil {
		t.Fatalf("❌ Failed to create integration server: %v", err)
	}
	t.Logf("🚀 Starting integration server...")
	integration, err := srv.Start()
	if err != nil {
		t.Fatalf("❌ Failed to create integration server: %v", err)
	}
	t.Cleanup(func() {
		t.Logf("🛑 Stopping integration server...")
		if err := integration.Stop(context.Background()); err != nil {
			t.Fatalf("❌ Failed to stop integration server: %v", err)
		}
	})

	baseURL := fmt.Sprintf("http://%s/api", integration.GetListen())
	var projects domain.SearchProjectsResponseBody
	getJSON(t, baseURL+"/projects?q=Search", &projects)
	if projects.TotalItems == nil || *projects.TotalItems == 0 {
		t.Fatal("the returned project list is empty")
	}

	var supplements *domain.ProjectSupplement
	getJSON(t, baseURL+"/projects/project-001/supplements", &supplements)
	if supplements == nil {
		t.Fatal("the returned project supplements are nil")
	}
}

func getJSON(t *testing.T, url string, target any) {
	t.Helper()

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET %s returned HTTP %d", url, response.StatusCode)
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode response from %s: %v", url, err)
	}
}
