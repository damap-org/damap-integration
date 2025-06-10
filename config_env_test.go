package damap_integration_test

import (
	"fmt"
	"testing"

	damap_integration "github.com/damap-org/damap-integration"
)

func TestEnvParsing(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       []string
		shouldError bool
		check       func(cfg damap_integration.Config) error
	}

	for name, tc := range map[string]testCase{
		"listen-addr-invalid": {
			input: []string{
				"DAMAP_INTEGRATION_LISTEN=0.0.0.",
			},
			shouldError: true,
		},
		"listen-addr": {
			input: []string{
				"DAMAP_INTEGRATION_LISTEN=0.0.0.0:8080",
			},
			check: func(cfg damap_integration.Config) error {
				if cfg.Listen == nil {
					return fmt.Errorf("incorrect listen address returned: nil")
				}
				if cfg.Listen.String() != "0.0.0.0:8080" {
					return fmt.Errorf("incorrect listen address returned: %s", cfg.Listen.String())
				}
				return nil
			},
		},
		"listen-addr-multiple": {
			input: []string{
				"DAMAP_INTEGRATION_LISTEN=0.0.0.0:8080",
				"DAMAP_INTEGRATION_LISTEN=0.0.0.0:8443",
			},
			check: func(cfg damap_integration.Config) error {
				if cfg.Listen == nil {
					return fmt.Errorf("incorrect listen address returned: nil")
				}
				if cfg.Listen.String() != "0.0.0.0:8443" {
					return fmt.Errorf("incorrect listen address returned: %s", cfg.Listen.String())
				}
				return nil
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Logf("⚙️ Setting up from %d environment variable(s)...", len(tc.input))
			c := damap_integration.Config{}
			if err := damap_integration.WithEnv(tc.input)(&c); err != nil {
				if tc.shouldError {
					t.Logf("✅ Environment setup returned correct error: %v", err)
					return
				}
				t.Fatalf("❌ Environment setup failed, unexpected error: %v", err)
			}
			if tc.shouldError {
				t.Fatalf("❌ Environment setup succeeded, expected error")
			}
			t.Logf("⚙️ Environment setup succeeded, checking configured values...")
			if err := tc.check(c); err != nil {
				t.Fatalf("❌ Invalid configured values: %v", err)
			}
			t.Logf("✅ Correct configuration values.")
		})
	}
}
