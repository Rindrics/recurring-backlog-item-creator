package main

import (
	"errors"
	"fmt"
)

func ValidateConfig(config Config) error {
	if config.Defaults.ProjectID == "" {
		return errors.New("defaults.project_id is required")
	}
	if config.Defaults.TargetRepo == "" {
		return errors.New("defaults.target_repo is required")
	}
	if len(config.Issues) == 0 {
		return errors.New("at least one issue is required")
	}

	for i, issue := range config.Issues {
		if err := ValidateIssue(issue); err != nil {
			return fmt.Errorf("issues[%d]: %w", i, err)
		}
	}

	return nil
}

func ValidateIssue(issue Issue) error {
	if issue.Name == "" {
		return errors.New("name is required")
	}
	if len(issue.CreationMonths) == 0 {
		return errors.New("creation_months is required and must not be empty")
	}
	if issue.TemplateFile == nil {
		return errors.New("template_file is required")
	}

	for i, month := range issue.CreationMonths {
		if month < 1 || month > 12 {
			return fmt.Errorf("creation_months[%d]: invalid month value %d (must be 1-12)", i, month)
		}
	}

	return nil
}
