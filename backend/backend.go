package backend

import (
	"github.com/damap-org/damap-integration/domain"
)

// Descriptor serves as a backend descriptor. Its only job is to return a config struct that can be filled based on
// the "cli" and "env" tags.
type Descriptor interface {
	// GetName returns a human-readable name for the backend.
	GetName() string
	// GetDescription returns a human-readable description for the backend.
	GetDescription() string
	// GetID returns this backend's unique ID.
	GetID() string
	// GetConfiguration returns a struct that can be configured using the "cli" and "env" tags.
	GetConfiguration() Configuration
}

// TypedBackendDescriptor is a specialized type so you can configure the backend manually if needed.
type TypedBackendDescriptor[T Configuration] interface {
	Descriptor

	// GetTypedConfiguration is identical to GetConfiguration, but returns the backend-specific type so you can
	// configure it manually.
	GetTypedConfiguration() T
}

// Configuration is a struct that has "cli" and "env" tags to fill it with configuration. The Build() function validates
// the values and builds a ConfiguredBackend.
type Configuration interface {
	// GetName returns a human-readable name for the backend.
	GetName() string
	// GetDescription returns a human-readable description for the backend.
	GetDescription() string
	// GetID returns this backend's unique ID.
	GetID() string
	// Configure validates the values in the struct and builds a ConfiguredBackend.
	Configure() (ConfiguredBackend, error)
}

// ConfiguredBackend is a fully configured backend that is ready to use.
type ConfiguredBackend interface {
	GetProjects() ([]domain.Project, error)
	GetPeople(projectId string) ([]domain.Person, error)
}
