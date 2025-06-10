package damap_integration_test

import (
	"testing"

	damap_integration "github.com/tuwien-csd/damap-integration"
)

func TestNewConfig(t *testing.T) {
	cfg := damap_integration.NewConfig()
	if cfg.Listen.String() != "127.0.0.1:8080" {
		t.Fatalf("invalid listen default value: %s", cfg.Listen.String())
	}
}
