package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

func outputJSON(ctx context.Context, issuesToCreate IssuesToCreate, defaults Defaults, ghClient GitHubClient) error {
	output := make([]IssueOutput, 0, len(issuesToCreate.Issues))

	// Track projects we've already logged
	loggedProjects := make(map[string]bool)

	for _, issue := range issuesToCreate.Issues {
		// Get target repo
		repo, err := issue.GetTargetRepo(defaults)
		if err != nil {
			return fmt.Errorf("failed to get target repo for issue %s: %w", issue.Name, err)
		}

		// Get project ID
		projectID := defaults.ProjectID
		if issue.ProjectID != nil {
			projectID = *issue.ProjectID
		}

		// Log project name if not already logged
		if !loggedProjects[projectID] {
			projectName, err := ghClient.GetProjectName(ctx, projectID)
			if err != nil {
				// Log error but continue
				log.Printf("Warning: failed to get project name for %s: %v", projectID, err)
			} else {
				log.Printf("Project: %s", projectName)
			}
			loggedProjects[projectID] = true
		}

		// Get project fields
		projectFields, err := ghClient.GetProjectFields(ctx, projectID, repo.Owner)
		if err != nil {
			return fmt.Errorf("failed to get project fields for issue %s: %w", issue.Name, err)
		}

		// Create field map for quick lookup
		fieldMap := make(map[string]ProjectField)
		for _, field := range projectFields {
			fieldMap[field.Name] = field
		}

		// Build field_updates array
		fieldUpdates := make([]FieldUpdate, 0, len(issue.Fields))
		for fieldName, fieldValue := range issue.Fields {
			field, exists := fieldMap[fieldName]
			if !exists {
				return fmt.Errorf("field '%s' not found in project for issue %s", fieldName, issue.Name)
			}

			fieldUpdate := FieldUpdate{
				FieldID:   field.ID,
				FieldType: field.DataType,
			}

			// Handle different field types
			switch field.DataType {
			case "TEXT", "NUMBER":
				fieldUpdate.Value = &fieldValue
			case "SINGLE_SELECT":
				// Find option ID by name
				var optionID *string
				for _, opt := range field.Options {
					if opt.Name == fieldValue {
						optionID = &opt.ID
						break
					}
				}
				if optionID == nil {
					return fmt.Errorf("option '%s' not found in field '%s' for issue %s", fieldValue, fieldName, issue.Name)
				}
				fieldUpdate.OptionID = optionID
			default:
				return fmt.Errorf("unsupported field type '%s' for field '%s' in issue %s", field.DataType, fieldName, issue.Name)
			}

			fieldUpdates = append(fieldUpdates, fieldUpdate)
		}

		// Generate title from name, title_prefix, and title_suffix
		expandedPrefix, err := expandTitlePrefix(issue.TitlePrefix)
		if err != nil {
			return fmt.Errorf("failed to expand title_prefix for issue %s: %w", issue.Name, err)
		}

		expandedSuffix, err := expandTitleSuffix(issue.TitleSuffix)
		if err != nil {
			return fmt.Errorf("failed to expand title_suffix for issue %s: %w", issue.Name, err)
		}

		// Build title: "{prefix} {name} {suffix}"
		// Add space between prefix and name only if prefix doesn't end with space
		// Add space between name and suffix only if suffix doesn't start with space
		title := issue.Name
		if expandedPrefix != "" {
			if strings.HasSuffix(expandedPrefix, " ") {
				title = expandedPrefix + title
			} else {
				title = expandedPrefix + " " + title
			}
		}
		if expandedSuffix != "" {
			if strings.HasPrefix(expandedSuffix, " ") {
				title = title + expandedSuffix
			} else {
				title = title + " " + expandedSuffix
			}
		}

		item := IssueOutput{
			Name:         issue.Name,
			Title:        title,
			TemplateFile: issue.TemplateFile,
			ProjectID:    issue.ProjectID,
			TargetRepo:   issue.TargetRepo,
			FieldUpdates: fieldUpdates,
		}

		output = append(output, item)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// expandTitleTemplate expands template variables in a title template string.
// If templateStr is nil or empty, returns an empty string.
// Supported template functions:
//   - {{Date}} - Current date in YYYY-MM-DD format
//   - {{Year}} - Current year (e.g., 2025)
//   - {{Month}} - Current month (e.g., 01)
//   - {{YearMonth}} - Current year and month in YYYY-MM format
func expandTitleTemplate(templateStr *string, templateName string) (string, error) {
	if templateStr == nil || *templateStr == "" {
		return "", nil
	}

	now := time.Now()
	funcMap := template.FuncMap{
		"Date": func() string {
			return now.Format("2006-01-02")
		},
		"Year": func() string {
			return now.Format("2006")
		},
		"Month": func() string {
			return now.Format("01")
		},
		"YearMonth": func() string {
			return now.Format("2006-01")
		},
	}

	tmpl, err := template.New("title").Funcs(funcMap).Parse(*templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s template: %w", templateName, err)
	}

	var buf bytes.Buffer
	// Execute with empty data since we're using functions, not data fields
	if err := tmpl.Execute(&buf, struct{}{}); err != nil {
		return "", fmt.Errorf("failed to execute %s template: %w", templateName, err)
	}

	return buf.String(), nil
}

// expandTitleSuffix expands template variables in title_suffix and returns the expanded suffix.
// If titleSuffix is nil or empty, returns an empty string.
func expandTitleSuffix(titleSuffix *string) (string, error) {
	return expandTitleTemplate(titleSuffix, "title_suffix")
}

// expandTitlePrefix expands template variables in title_prefix and returns the expanded prefix.
// If titlePrefix is nil or empty, returns an empty string.
func expandTitlePrefix(titlePrefix *string) (string, error) {
	return expandTitleTemplate(titlePrefix, "title_prefix")
}
