package damap_integration

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/damap-org/damap-integration/backend"
	"github.com/damap-org/damap-integration/server"
)

func New(
	opts ...Opt,
) (Integration, error) {
	cfg := NewConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	var keyPair tls.Certificate
	if cfg.Certificate != "" && cfg.Key != "" {
		cert := cfg.Certificate
		key := cfg.Key
		var err error
		if !strings.HasPrefix(strings.TrimSpace(cert), "-----") {
			certData, err := os.ReadFile(cert)
			if err != nil {
				return nil, fmt.Errorf("failed to read certificate file %s: %w", cfg.Certificate, err)
			}
			cert = string(certData)
		}
		if !strings.HasPrefix(strings.TrimSpace(key), "-----") {
			keyData, err := os.ReadFile(key)
			if err != nil {
				return nil, fmt.Errorf("failed to read key file %s: %w", cfg.Key, err)
			}
			cert = string(keyData)
		}

		keyPair, err = tls.X509KeyPair([]byte(cert), []byte(key))
		if err != nil {
			return nil, fmt.Errorf("failed to load PEM-encoded certificate and key: %w", err)
		}
	}

	configuredBackend, err := cfg.backendConfigurations[cfg.Backend].Configure()
	if err != nil {
		return nil, fmt.Errorf("invalid backend configuration: %w", err)
	}
	return &integration{
		configuredBackend,
		*cfg,
		keyPair,
	}, nil
}

type integration struct {
	backend backend.ConfiguredBackend
	config  Config
	keyPair tls.Certificate
}

func (i integration) Start() (RunningIntegration, error) {
	logger := i.config.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	result := &runningIntegration{
		backend: i.backend,
	}

	var listeners []net.Listener
	if i.config.Listen != nil {
		ln, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(*i.config.Listen))
		if err != nil {
			logger.Error("Failed to listen on TCP socket", "address", *i.config.Listen, "cause", err)
			return nil, fmt.Errorf("failed to listen on %s (%w)", i.config.Listen.String(), err)
		}
		listeners = append(listeners, ln)
		addrPort := ln.Addr().(*net.TCPAddr).AddrPort()
		result.listen = &addrPort
	}
	if i.config.ListenTLS != nil {
		ln, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(*i.config.Listen))
		if err != nil {
			for _, ln := range listeners {
				_ = ln.Close()
			}
			logger.Error("Failed to listen on TCP socket", "address", *i.config.ListenTLS, "cause", err)
			return nil, fmt.Errorf("failed to listen on %s (%w)", i.config.ListenTLS.String(), err)
		}
		tlsListener := tls.NewListener(ln, &tls.Config{
			Certificates: []tls.Certificate{
				i.keyPair,
			},
			MinVersion: tls.VersionTLS13,
			NextProtos: []string{"h2", "http/1.1"},
		})
		listeners = append(listeners, tlsListener)
		addrPort := ln.Addr().(*net.TCPAddr).AddrPort()
		result.tlsListen = &addrPort
	}

	result.server = &http.Server{Handler: server.NewHandler(i.backend), ReadHeaderTimeout: 5 * time.Second}
	var once sync.Once
	result.done = make(chan struct{})
	for _, ln := range listeners {
		go func() {
			logger.Info("DAMAP integration now available", "address", ln.Addr())
			if err := result.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				result.lastError = err
			}
			once.Do(func() {
				logger.Info("DAMAP integration shutting down")
				close(result.done)
				_ = result.server.Shutdown(context.Background())
			})
		}()
	}

	return result, nil
}

type runningIntegration struct {
	backend   backend.ConfiguredBackend
	server    *http.Server
	lastError error
	done      chan struct{}
	listen    *netip.AddrPort
	tlsListen *netip.AddrPort
}

func (r runningIntegration) GetListen() *netip.AddrPort {
	return r.listen
}

func (r runningIntegration) GetTLSListen() *netip.AddrPort {
	return r.tlsListen
}

func (r runningIntegration) OnStop() chan struct{} {
	return r.done
}

func (r runningIntegration) Stop(ctx context.Context) error {
	if err := r.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return r.lastError
}
