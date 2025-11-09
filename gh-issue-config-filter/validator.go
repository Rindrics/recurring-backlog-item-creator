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

	// Get project name for error messages
	ctx := context.Background()
	projectName, err := ghClient.GetProjectName(ctx, projectID)
	if err != nil {
		// Fallback to project ID if name cannot be retrieved
		Debugf("Failed to get project name for %s: %v, using project ID as fallback", projectID, err)
		projectName = projectID
	} else if projectName == "" {
		// If project name is empty, use project ID as fallback
		Debugf("Project name is empty for %s, using project ID as fallback", projectID)
		projectName = projectID
	} else {
		Debugf("Successfully retrieved project name: %s for project ID: %s", projectName, projectID)
	}

	// Get project fields for validation
	projectFields, err := ghClient.GetProjectFields(ctx, projectID, issueRepo.Owner)
	if err != nil {
		return fmt.Errorf("failed to get project fields: %w", err)
	}

	// Create a map of field names to fields for quick lookup
	fieldMap := make(map[string]ProjectField)
	for _, field := range projectFields {
		fieldMap[field.Name] = field
		Debugf("Found project field: %s (ID: %s, Type: %s)", field.Name, field.ID, field.DataType)
	}

	// Validate fields
	// Format project display name: "Name (ID)" if name is different from ID, otherwise just ID
	projectDisplayName := projectID
	if projectName != "" && projectName != projectID {
		projectDisplayName = fmt.Sprintf("%s (%s)", projectName, projectID)
	}
	Debugf("Using project display name: %s", projectDisplayName)
	if err := ValidateIssueFields(issue, fieldMap, projectDisplayName); err != nil {
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
		if !month.IsValid() {
			return fmt.Errorf("creation_months[%d]: invalid month value %d (must be 1-12)", i, month)
		}
	}

	return nil
}

func ValidateIssueFields(issue Issue, fieldMap map[string]ProjectField, projectName string) error {
	for fieldName, fieldValue := range issue.Fields {
		Debugf("Validating field '%s' with value '%s'", fieldName, fieldValue)
		field, exists := fieldMap[fieldName]
		if !exists {
			// Get list of available field names
			availableFields := make([]string, 0, len(fieldMap))
			for name := range fieldMap {
				availableFields = append(availableFields, name)
			}
			return fmt.Errorf("field '%s' does not exist in project '%s'. Available fields: %v", fieldName, projectName, availableFields)
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
