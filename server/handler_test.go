package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/damap-org/damap-integration/backend"
	"github.com/damap-org/damap-integration/domain"
	"github.com/damap-org/damap-integration/server"
)

type noLeaderBackend struct {
	backend.ConfiguredBackend
}

func (noLeaderBackend) GetLeaderForProject(id string) (*domain.Person, error) {
	return nil, fmt.Errorf("%w: project %q has no marked leader", domain.ErrNotFound, id)
}

func TestGetLeaderForProjectNotFound(t *testing.T) {
	response := httptest.NewRecorder()
	server.NewHandler(noLeaderBackend{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/projects/42/leader", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected HTTP 404, got %d: %s", response.Code, response.Body.String())
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected application/json, got %q", got)
	}
	var body domain.NotFoundError
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if want := `not found: project "42" has no marked leader`; body.Message != want {
		t.Fatalf("expected message %q, got %q", want, body.Message)
	}
}
