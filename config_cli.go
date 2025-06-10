package damap_integration

import (
	"fmt"
	"reflect"
	"strings"
)

// WithCLI sets all configuration options from command-line parameters.
// The parameters should be in the format of --option=value or --option value (separate entries in the list)
// in order to be parsed properly. The first element should always be the program name and will be ignored.
// You can pass os.Args directly to this function.
func WithCLI(args []string) func(config *Config) error {
	return func(config *Config) error {
		if len(args) == 0 {
			return fmt.Errorf("no arguments provided")
		}
		cliTags := map[string]reflect.Value{}
		cliTagFields := map[string]reflect.StructField{}
		configValue := reflect.ValueOf(config).Elem()
		for i := 0; i < configValue.NumField(); i++ {
			cliTag := configValue.Type().Field(i).Tag.Get("cli")
			if cliTag != "" {
				cliTags[cliTag] = configValue.Field(i)
				cliTagFields[cliTag] = configValue.Type().Field(i)
			}
		}

		for id, b := range config.backendConfigurations {
			configVal := reflect.ValueOf(b)
			if configVal.Type().Kind() == reflect.Ptr {
				configVal = configVal.Elem()
			}
			for i := 0; i < configVal.NumField(); i++ {
				cliTag := configVal.Type().Field(i).Tag.Get("cli")
				if cliTag != "" {
					cliTag = transformCLIOption(cliTag, id)
					cliTags[cliTag] = configVal.Field(i)
					cliTagFields[cliTag] = configVal.Type().Field(i)
				}
			}
		}

		// Note: this could be implemented with Viper or Cobra, but that would be less explicit and introduce an
		// unnecessary dependency for something relatively trivial.
		args = args[1:]
		optValues := map[string]string{}
		for i := 0; i < len(args); i++ {
			if !strings.HasPrefix(args[i], "--") {
				return fmt.Errorf("unexpected argument: %s", args[i])
			}
			key := args[i]
			var value string
			parts := strings.SplitN(args[i], "=", 2)
			if len(parts) == 2 {
				key = parts[0]
				value = parts[1]
			} else {
				if i == len(args)-1 || strings.HasPrefix(args[i+1], "--") {
					return fmt.Errorf("%s requires an argument", args[i])
				}
				value = args[i+1]
				i++
			}
			opt := strings.TrimPrefix(key, "--")
			if _, ok := optValues[opt]; ok {
				return fmt.Errorf("duplicate option: %s", args[i])
			}
			if _, ok := cliTags[opt]; !ok {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			optValues[opt] = value
		}

		for opt, val := range optValues {
			value := reflect.ValueOf(val)
			targetField := cliTagFields[opt]
			targetFieldValue := cliTags[opt]
			if err := setStructFieldValue(targetFieldValue, value, targetField); err != nil {
				return fmt.Errorf("invalid value for option --%s: %s (%w)", opt, val, err)
			}
		}
		return nil
	}
}

func transformCLIOption(cliOption string, id string) string {
	return id + "-backend-" + cliOption
}
