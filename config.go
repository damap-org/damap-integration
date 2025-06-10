package damap_integration

import (
	"fmt"
	"log/slog"
	"net/netip"
	"reflect"

	"github.com/damap-org/damap-integration/backend"
)

// NewConfig returns a Config preconfigured with the default values.
func NewConfig() *Config {
	cfg := &Config{}
	cfg.backendConfigurations = make(map[string]backend.Configuration)
	for id, b := range backend.Backends {
		cfg.backendConfigurations[id] = b.GetConfiguration()
		reflected := reflect.ValueOf(cfg.backendConfigurations[id])
		if reflected.Type().Kind() != reflect.Ptr && reflected.Type().Kind() != reflect.Struct {
			panic(fmt.Errorf(
				"invalid backend %s, expected a struct or pointer to struct as a configuration, got %v",
				cfg.backendConfigurations[id].GetID(),
				reflected.Type().String(),
			))
		}
	}
	reflectedCfg := reflect.ValueOf(cfg).Elem()
	for i := 0; i < reflectedCfg.NumField(); i++ {
		defaultValue := reflectedCfg.Type().Field(i).Tag.Get("default")
		if defaultValue == "" {
			continue
		}
		targetFieldValue := reflectedCfg.Field(i)
		if err := setStructFieldValue(targetFieldValue, reflect.ValueOf(defaultValue), reflectedCfg.Type().Field(i)); err != nil {
			panic(fmt.Errorf("bug: cannot set the default value for config field %s (%w)", reflectedCfg.Type().Field(i).Name, err))
		}
	}
	return cfg
}

type Config struct {
	// Logger is the sink to send the logs to.
	Logger *slog.Logger

	// Listen holds the IP address and port to bind the non-secure listener to.
	// Set this to "" to disable listening on a non-secure port.
	Listen *netip.AddrPort `default:"127.0.0.1:8080" cli:"listen" env:"DAMAP_INTEGRATION_LISTEN" help:"Address to listen on without encryption."`

	// ListenTLS holds the IP address and port to bind the secure listener to.
	// Set this to "" to disable listening on a secure port.
	ListenTLS *netip.AddrPort `default:"" cli:"listen-tls" env:"DAMAP_INTEGRATION_LISTEN_TLS" help:"Address to listen on with encryption. Requires the --certificate and --key options."`

	// Certificate holds the PEM-encoded certificate or a reference to a file containing the PEM-encoded
	// certificate to use for ListenTLS.
	Certificate string `default:"" cli:"certificate" env:"DAMAP_INTEGRATION_CERTIFICATE" help:"PEM-encoded x509 certificate or file containing the PEM-encoded certificate."`

	// Key holds the PEM-encoded private key or a reference to a file containing a PEM-encoded
	// private key to use for ListenTLS.
	Key string `default:"" cli:"key" env:"DAMAP_INTEGRATION_KEY" help:"PEM-encoded private key or file containing the PEM-encoded private key."`

	// Backend selects which backend to use.
	Backend string `default:"dummy" cli:"backend" env:"DAMAP_INTEGRATION_BACKEND" help:"Backend to use (use backend ID in this field)."`

	backendConfigurations map[string]backend.Configuration
}

// ConfigHelp returns a help text for the CLI and env configurations.
func ConfigHelp() string {
	result := "Usage: damap-integration [options]\n\n"
	result += "Options:\n"
	result += "  -h, --help: show this help message.\n"
	t := reflect.ValueOf(Config{}).Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		cli := f.Tag.Get("cli")
		env := f.Tag.Get("env")
		def := f.Tag.Get("default")
		help := f.Tag.Get("help")
		if cli == "" && env == "" {
			continue
		}
		if cli != "" {
			result += "  --" + cli + "=value\n"
		}
		if env != "" {
			result += "    Environment variable: " + env + "=value\n"
		}
		if def != "" {
			result += "    Default: " + def + "\n"
		}
		if help != "" {
			result += "    " + help + "\n"
		}
		result += "\n"
	}
	result += "\n"
	for _, b := range backend.Backends {
		title := b.GetName() + " backend (ID: " + b.GetID() + ")"
		result += title + "\n"
		for range len(title) {
			result += "="
		}
		result += "\n\n"
		if description := b.GetDescription(); description != "" {
			result += description + "\n\n"
		}
		t := reflect.ValueOf(b.GetConfiguration()).Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		optCount := 0
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			cli := f.Tag.Get("cli")
			env := f.Tag.Get("env")
			def := f.Tag.Get("default")
			help := f.Tag.Get("help")
			if cli == "" && env == "" {
				continue
			}
			if cli != "" {
				result += "  --" + transformCLIOption(cli, b.GetID()) + "=value\n"
			}
			if env != "" {
				result += "    Environment variable: " + transformEnvVariable(env, b.GetID()) + "=value\n"
			}
			if def != "" {
				result += "    Default: " + def + "\n"
			}
			if help != "" {
				result += "    " + help + "\n"
			}
			optCount++
		}
		if optCount == 0 {
			result += "No configuration options for this backend.\n"
		}
		result += "\n"
	}
	return result
}

func (c *Config) Validate() error {
	if c.Listen == nil && c.ListenTLS == nil {
		return fmt.Errorf("both the plain text and encrypted listen socket are disabled; please provide either listen or listen-tls")
	}
	if c.ListenTLS != nil && (c.Certificate == "" || c.Key == "") {
		return fmt.Errorf("no certificate or key specified for the encrypted listen socket; please specify --certificate and --key")
	}
	return nil
}

// Opt is a function that sets up a specific aspect of Config.
type Opt func(config *Config) error

// WithLogger adds a logger to the integration.
func WithLogger(logger *slog.Logger) Opt {
	return func(config *Config) error {
		config.Logger = logger
		return nil
	}
}

// WithListen registers the plain text listen address the integration server should run on. Pass an empty string to
// disable the plain text port.
func WithListen(addr string) func(config *Config) error {
	return func(config *Config) error {
		if addr == "" {
			config.Listen = nil
			return nil
		}
		addrPort, err := netip.ParseAddrPort(addr)
		if err != nil {
			return fmt.Errorf("invalid listen address: %s (%w)", addr, err)
		}
		config.Listen = &addrPort
		return nil
	}
}

// WithListenTLS sets up the encrypted address the integration server should run on. Pass an empty string to disable the
// TLS port.
func WithListenTLS(addr string) func(config *Config) error {
	return func(config *Config) error {
		if addr == "" {
			config.Listen = nil
			return nil
		}
		addrPort, err := netip.ParseAddrPort(addr)
		if err != nil {
			return fmt.Errorf("invalid listen-tls address: %s (%w)", addr, err)
		}
		config.ListenTLS = &addrPort
		return nil
	}
}

// WithCertificate configures a PEM-encoded x509 certificate to be used as a server certificate.
// You can pass either a file or the PEM-encoded certificate directly.
func WithCertificate(certificate string) func(config *Config) error {
	return func(config *Config) error {
		config.Certificate = certificate
		return nil
	}
}

// WithKey configures a PEM-encoded private key to be used for connection encryption.
// You can pass either a file or a PEM-encoded certificate directly.
func WithKey(key string) func(config *Config) error {
	return func(config *Config) error {
		config.Key = key
		return nil
	}
}

// WithBackend passes a backend configuration that can be used to obtain data. This backend can be a backend that is
// not registered in backend.Backends.
func WithBackend(backendConfiguration backend.Configuration) func(config *Config) error {
	return func(config *Config) error {
		backendId := backendConfiguration.GetID()
		config.Backend = backendId
		config.backendConfigurations[backendId] = backendConfiguration
		return nil
	}
}

func setStructFieldValue(targetFieldValue reflect.Value, value reflect.Value, targetField reflect.StructField) error {
	targetType := targetFieldValue.Type()
	if !value.Type().AssignableTo(targetType) {
		if targetType.Kind() == reflect.Ptr && targetFieldValue.IsNil() {
			// If the value is currently nil, we need to make it a non-nil value so we can actually
			// call the UnmarshalText function on it.
			targetFieldValue.Set(reflect.New(targetType.Elem()))
		}
		// Get the receiver method so we can call it
		unmarshalTextFunc := targetFieldValue.MethodByName("UnmarshalText")
		if !unmarshalTextFunc.IsValid() {
			panic(fmt.Errorf(
				"bug: the Config field %s is not assignable from a string and does not have an UnmarshalText function either",
				targetField.Name,
			))
		}
		returnValues := unmarshalTextFunc.Call([]reflect.Value{reflect.ValueOf([]byte(value.String()))})
		if len(returnValues) != 1 {
			panic(fmt.Errorf(
				"bug: the Config field %s has an UnmarshalText function, but it does not have a single error return value",
				targetField.Name,
			))
		}
		if !returnValues[0].IsNil() {
			err, ok := returnValues[0].Interface().(error)
			if !ok {
				panic(fmt.Errorf(
					"bug: the Config field %s has an UnmarshalText function, but it does not have a single error return value",
					targetField.Name,
				))
			}
			return err
		}
	} else {
		targetFieldValue.Set(value)
	}
	return nil
}
