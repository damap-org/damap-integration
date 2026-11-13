package client_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/damap-org/damap-integration/client"
	"github.com/damap-org/damap-integration/domain"
)

func TestClientNotFound(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		call func(context.Context, client.Client) error
	}{
		{"project", "/api/projects/42", func(ctx context.Context, c client.Client) error {
			_, err := c.GetProject(ctx, "42")
			return err
		}},
		{"person", "/api/persons/42", func(ctx context.Context, c client.Client) error {
			_, err := c.GetPerson(ctx, "42")
			return err
		}},
		{"leader", "/api/projects/42/leader", func(ctx context.Context, c client.Client) error {
			_, err := c.GetLeaderForProject(ctx, "42")
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, withJSON := range []bool{true, false} {
				name := "without JSON"
				if withJSON {
					name = "with JSON"
				}
				t.Run(name, func(t *testing.T) {
					srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path != tc.path {
							t.Errorf("expected path %q, got %q", tc.path, r.URL.Path)
						}
						if withJSON {
							w.Header().Set("Content-Type", "application/json")
						}
						w.WriteHeader(http.StatusNotFound)
						if withJSON {
							_, _ = w.Write([]byte(`{"message":"lookup returned no result"}`))
						}
					}))
					defer srv.Close()

					c, err := client.New(srv.URL + "/api/")
					if err != nil {
						t.Fatal(err)
					}
					err = tc.call(t.Context(), c)
					if !errors.Is(err, domain.ErrNotFound) {
						t.Fatalf("expected ErrNotFound, got %v", err)
					}
					if withJSON && !strings.Contains(err.Error(), "lookup returned no result") {
						t.Fatalf("expected server message, got %v", err)
					}
				})
			}
		})
	}
}
