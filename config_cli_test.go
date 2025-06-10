package damap_integration_test

import (
	"fmt"
	"testing"

	"github.com/damap-org/damap-integration"
)

func TestWithCLI(t *testing.T) {
	type testCase struct {
		input       []string
		shouldError bool
		check       func(cfg damap_integration.Config) error
	}

	for name, tc := range map[string]testCase{
		"listen-addr-invalid": {
			input: []string{
				"--listen=0.0.0.",
			},
			shouldError: true,
		},
		"listen-addr-invalid-split": {
			input: []string{
				"--listen", "0.0.0.",
			},
			shouldError: true,
		},
		"listen-addr-missing": {
			input: []string{
				"--listen",
			},
			shouldError: true,
		},
		"listen-addr": {
			input: []string{
				"--listen=0.0.0.0:8080",
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
		"listen-addr-split": {
			input: []string{
				"--listen", "0.0.0.0:8080",
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
				"--listen=0.0.0.0:8080",
				"--listen=0.0.0.0:8443",
			},
			shouldError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Logf("⚙️ Setting up from %d cli option(s)...", len(tc.input))
			c := damap_integration.Config{}
			if err := damap_integration.WithCLI(append([]string{"foo"}, tc.input...))(&c); err != nil {
				if tc.shouldError {
					t.Logf("✅ CLI setup returned correct error: %v", err)
					return
				}
				t.Fatalf("❌ CLI setup failed, unexpected error: %v", err)
			}
			if tc.shouldError {
				t.Fatalf("❌ CLI setup succeeded, expected error")
			}
			t.Logf("⚙️ CLI setup succeeded, checking configured values...")
			if err := tc.check(c); err != nil {
				t.Fatalf("❌ Invalid configured values: %v", err)
			}
			t.Logf("✅ Correct configuration values.")
		})
	}
}
