package dummy

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/tuwien-csd/damap-integration/backend"
	"github.com/tuwien-csd/damap-integration/domain"
)

//go:embed people.json
var peopleData []byte

//go:embed projects.json
var projectsData []byte

const ID = "dummy"

func init() {
	var people []domain.Person
	var projects []domain.Project
	if err := json.Unmarshal(peopleData, &people); err != nil {
		panic(fmt.Errorf("bug: failed to unmarshal dummy people data: %w", err))
	}
	if err := json.Unmarshal(projectsData, &projects); err != nil {
		panic(fmt.Errorf("bug: failed to unmarshal dummy projects data: %w", err))
	}
	backend.Backends[ID] = &dummyBackend{
		people,
		projects,
	}
}

// New returns the Dummy backend.
func New() backend.Descriptor {
	// Backends are stateless, just return
	return backend.Backends[ID]
}

type dummyBackend struct {
	people   []domain.Person
	projects []domain.Project
}

func (d *dummyBackend) GetName() string {
	return "Dummy"
}

func (d *dummyBackend) GetDescription() string {
	return "A backend that returns dummy data."
}

func (d *dummyBackend) Configure() (backend.ConfiguredBackend, error) {
	return d, nil
}

func (d *dummyBackend) GetID() string {
	return ID
}

func (d *dummyBackend) GetConfiguration() backend.Configuration {
	return d
}

func (d *dummyBackend) GetProjects() ([]domain.Project, error) {
	return d.projects, nil
}

func (d *dummyBackend) GetPeople(_ string) ([]domain.Person, error) {
	return d.people, nil
}
