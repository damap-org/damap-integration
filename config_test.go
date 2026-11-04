package damap_integration_test

import (
	"testing"

	damap_integration "github.com/damap-org/damap-integration"
)

func TestNewConfig(t *testing.T) {
	cfg := damap_integration.NewConfig()
	if cfg.Listen.String() != "127.0.0.1:8090" {
		t.Fatalf("invalid listen default value: %s", cfg.Listen.String())
	}
}
