package damap_integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/damap-org/damap-integration"
	"github.com/damap-org/damap-integration/client"
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

	t.Logf("⚙️ Creating integration client...")
	//goland:noinspection HttpUrlsUsage
	cli, err := client.New(fmt.Sprintf("http://%s/api/", integration.GetListen()))
	if err != nil {
		t.Fatalf("❌ Failed to create client: %v", err)
	}
	t.Logf("📂 Requesting project list...")
	projects, err := cli.SearchProjects(context.Background(), &domain.SearchProjectsParams{Query: "Search"})
	if err != nil {
		t.Fatalf("❌ Failed to request project list: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("❌ The returned project list is empty!")
	}

	t.Logf("📂 Requesting project supplements...")
	supplements, err := cli.GetProjectSupplements(context.Background(), "project-001")
	if err != nil {
		t.Fatalf("❌ Failed to request project supplements: %v", err)
	}
	if supplements == nil {
		t.Fatal("❌ The returned project supplements are nil")
	}
}
