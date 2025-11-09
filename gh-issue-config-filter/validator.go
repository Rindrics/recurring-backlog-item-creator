package main

import (
	"context"
	"errors"
	"fmt"
)

func ValidateConfig(config Config, ghClient GitHubClient) error {
	if config.Defaults.ProjectID == "" {
		return errors.New("defaults.project_id is required")
	}
	if config.Defaults.TargetRepo == "" {
		return errors.New("defaults.target_repo is required")
	}
	if len(config.Issues) == 0 {
		return errors.New("at least one issue is required")
	}

	// Validate each issue
	for i, issue := range config.Issues {
		if err := ValidateIssueWithProject(issue, config.Defaults, ghClient); err != nil {
			return fmt.Errorf("issues[%d]: %w", i, err)
		}
	}

	return nil
}

func ValidateIssueWithProject(issue Issue, defaults Defaults, ghClient GitHubClient) error {
	// Basic issue validation
	if err := ValidateIssue(issue); err != nil {
		return err
	}

	// Validate target_repo format
	issueRepo, err := issue.GetTargetRepo(defaults)
	if err != nil {
		return fmt.Errorf("invalid target_repo: %w", err)
	}

	// Get project ID for this issue
	projectID := defaults.ProjectID
	if issue.ProjectID != nil {
		projectID = *issue.ProjectID
	}

	// Get project fields for validation
	ctx := context.Background()
	projectFields, err := ghClient.GetProjectFields(ctx, projectID, issueRepo.Owner)
	if err != nil {
		return fmt.Errorf("failed to get project fields: %w", err)
	}

	// Create a map of field names to fields for quick lookup
	fieldMap := make(map[string]ProjectField)
	for _, field := range projectFields {
		fieldMap[field.Name] = field
	}

	// Validate fields
	if err := ValidateIssueFields(issue, fieldMap); err != nil {
		return err
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

func ValidateIssueFields(issue Issue, fieldMap map[string]ProjectField) error {
	for fieldName, fieldValue := range issue.Fields {
		field, exists := fieldMap[fieldName]
		if !exists {
			return fmt.Errorf("field '%s' does not exist in project", fieldName)
		}

		// For single-select fields, validate that the option exists
		if field.DataType == "SINGLE_SELECT" {
			optionExists := false
			for _, option := range field.Options {
				if option.Name == fieldValue {
					optionExists = true
					break
				}
			}
			if !optionExists {
				return fmt.Errorf("field '%s': option '%s' does not exist", fieldName, fieldValue)
			}
		}
	}

	return nil
}
