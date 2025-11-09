package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

const defaultConfigFile = ".gh-issue-config-filter.yml"

func main() {
	var (
		month      = flag.Int("month", 0, "Month (1-12) to filter issues")
		configFile = flag.String("config", "", "Path to config file")
		debug      = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	SetDebugMode(*debug)

	monthEnum, err := ParseMonth(*month)
	if err != nil {
		log.Fatalf("failed to parse month: %v", err)
	}

	configPath := *configFile
	if configPath == "" {
		configPath = defaultConfigFile
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("config file not found: %s\nUse --config to specify a different config file", configPath)
		}
		log.Fatalf("failed to load config: %v", err)
	}

	if err := ValidateConfig(config); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	issuesToCreate := GetIssuesToCreate(config, monthEnum)

	if err := outputJSON(issuesToCreate); err != nil {
		log.Fatalf("failed to output JSON: %v", err)
	}
}

func outputJSON(issuesToCreate IssuesToCreate) error {
	output := make([]map[string]interface{}, 0, len(issuesToCreate.Issues))
	for _, issue := range issuesToCreate.Issues {
		item := map[string]interface{}{
			"name":          issue.Name,
			"template_file": issue.TemplateFile,
			"fields":        issue.Fields,
			"project_id":    issue.ProjectID,
			"target_repo":   issue.TargetRepo,
		}
		if issue.TitleSuffix != nil {
			item["title_suffix"] = issue.TitleSuffix
		}
		output = append(output, item)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
