package damap_integration

import (
	"fmt"
	"reflect"
	"strings"
)

// WithEnv sets all configuration options from the passed OS environment.
// The environment entries should be in the form of key=value, the same as how
// os.Environ() returns values.
func WithEnv(env []string) func(config *Config) error {
	return func(config *Config) error {
		envTags := map[string]reflect.Value{}
		envTagFields := map[string]reflect.StructField{}
		configValue := reflect.ValueOf(config).Elem()
		for i := 0; i < configValue.NumField(); i++ {
			envTag := configValue.Type().Field(i).Tag.Get("env")
			if envTag != "" {
				envTags[envTag] = configValue.Field(i)
				envTagFields[envTag] = configValue.Type().Field(i)
			}
		}
		for id, b := range config.backendConfigurations {
			configVal := reflect.ValueOf(b)
			if configVal.Type().Kind() == reflect.Ptr {
				configVal = configVal.Elem()
			}
			for i := 0; i < configVal.NumField(); i++ {
				envTag := configVal.Type().Field(i).Tag.Get("env")
				if envTag != "" {
					envTag = transformEnvVariable(envTag, id)
					envTags[envTag] = configVal.Field(i)
					envTagFields[envTag] = configVal.Type().Field(i)
				}
			}
		}

		for _, item := range env {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) != 2 {
				continue
			}
			if _, ok := envTags[parts[0]]; !ok {
				continue
			}
			value := reflect.ValueOf(parts[1])
			targetField := envTagFields[parts[0]]
			targetFieldValue := envTags[parts[0]]
			if err := setStructFieldValue(targetFieldValue, value, targetField); err != nil {
				return fmt.Errorf("invalid value for environment variable %s: %s (%w)", parts[0], parts[1], err)
			}
		}
		return nil
	}
}

func transformEnvVariable(envVariable string, id string) string {
	return "DAMAP_INTEGRATION_" + strings.ToUpper(id) + "_BACKEND_" + strings.TrimPrefix(envVariable, "DAMAP_INTEGRATION_")
}
