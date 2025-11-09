package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configFile string) (Config, error) {
	Debug("loading config file: ", configFile)

	var config = Config{}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// Provide a more user-friendly error message for YAML parsing errors
		if strings.Contains(err.Error(), "unmarshal") {
			schema, schemaErr := generateConfigSchema()
			if schemaErr != nil {
				return config, fmt.Errorf("invalid YAML format in config file: %w", err)
			}
			return config, fmt.Errorf("invalid YAML format in config file: %w\n\nExpected schema:\n%s", err, schema)
		}
		return config, err
	}

	Debug("loaded config file: ", &config)

	return config, nil
}

// generateConfigSchema generates a YAML schema example from the Config struct
func generateConfigSchema() (string, error) {
	exampleConfig := Config{
		Defaults: Defaults{
			ProjectID:  "PVT_xxx",
			TargetRepo: "owner/repo",
		},
		Issues: []Issue{
			{
				Name:           "Example Issue",
				CreationMonths: []Month{January, February},
				TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/example.md"),
				Fields: map[string]string{
					"SP":     "5",
					"status": "Ready",
				},
			},
		},
	}

	data, err := yaml.Marshal(&exampleConfig)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ParseMonth(digit int) (Month, error) {
	if digit < 1 || digit > 12 {
		return 0, errors.New("month must be between 1 and 12")
	}
	return Month(digit), nil
}
