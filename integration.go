package damap_integration

import (
	"context"
	"net/netip"
)

// Integration is a configured integration that can be started using the Start() function.
type Integration interface {
	// Start starts the server and returns immediately.
	Start() (RunningIntegration, error)
}

// RunningIntegration is the Integration that has been started.
type RunningIntegration interface {
	// GetListen returns the listen IP and port if configured. This function is useful if you set the port to
	// 0 to allocate a dynamic port.
	GetListen() *netip.AddrPort
	// GetTLSListen returns the TLS listen IP and port if configured. This function is useful if you set the port to
	// 0 to allocate a dynamic port.
	GetTLSListen() *netip.AddrPort

	// OnStop returns a channel that closes when the server stops.
	OnStop() chan struct{}
	// Stop requests that the server stops. The passed context provides a timeout
	// for the shutdown. When the context expires, the existing connections are
	// aborted.
	Stop(ctx context.Context) error
}
