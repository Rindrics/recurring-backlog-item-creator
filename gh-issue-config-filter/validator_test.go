package main

import (
	"context"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	cases := []struct {
		name                string
		config              Config
		mockFields          []ProjectField
		expectError         bool
		expectErrorContains string
	}{
		{
			name: "valid config",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						Fields:         map[string]string{},
					},
				},
			},
			mockFields:          []ProjectField{},
			expectError:         false,
			expectErrorContains: "",
		},
		{
			name: "invalid - empty project_id",
			config: Config{
				Defaults: Defaults{
					ProjectID: "",
				},
			},
			mockFields:          []ProjectField{},
			expectError:         true,
			expectErrorContains: "defaults.project_id is required",
		},
		{
			name: "invalid - empty target_repo",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "",
				},
			},
			mockFields:          []ProjectField{},
			expectError:         true,
			expectErrorContains: "defaults.target_repo is required",
		},
		{
			name: "invalid - no issues",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{},
			},
			mockFields:          []ProjectField{},
			expectError:         true,
			expectErrorContains: "at least one issue is required",
		},
		{
			name: "invalid - issue with empty name",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{{Name: ""}},
			},
			mockFields:          []ProjectField{},
			expectError:         true,
			expectErrorContains: "name is required",
		},
		{
			name: "invalid - field does not exist",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						Fields: map[string]string{
							"NonExistentField": "value",
						},
					},
				},
			},
			mockFields: []ProjectField{
				{ID: "PVTFL_1", Name: "Status", DataType: "SINGLE_SELECT"},
			},
			expectError:         true,
			expectErrorContains: "field 'NonExistentField' does not exist in project",
		},
		{
			name: "invalid - single-select option does not exist",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						Fields: map[string]string{
							"Status": "InvalidOption",
						},
					},
				},
			},
			mockFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Status",
					DataType: "SINGLE_SELECT",
					Options: []ProjectFieldOption{
						{ID: "OPT_1", Name: "Ready"},
						{ID: "OPT_2", Name: "In Progress"},
					},
				},
			},
			expectError:         true,
			expectErrorContains: "field 'Status': option 'InvalidOption' does not exist",
		},
		{
			name: "valid - field exists",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						Fields: map[string]string{
							"Status": "Ready",
						},
					},
				},
			},
			mockFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Status",
					DataType: "SINGLE_SELECT",
					Options: []ProjectFieldOption{
						{ID: "OPT_1", Name: "Ready"},
						{ID: "OPT_2", Name: "In Progress"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid - field does not exist in default project",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						Fields: map[string]string{
							"NonExistentField": "value",
						},
					},
				},
			},
			mockFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Status",
					DataType: "SINGLE_SELECT",
					Options: []ProjectFieldOption{
						{ID: "OPT_1", Name: "Ready"},
					},
				},
			},
			expectError:         true,
			expectErrorContains: "field 'NonExistentField' does not exist in project",
		},
		{
			name: "invalid - field does not exist in overridden project",
			config: Config{
				Defaults: Defaults{
					ProjectID:  "default_project_id",
					TargetRepo: "default/repo",
				},
				Issues: []Issue{
					{
						Name:           "test",
						CreationMonths: []Month{January},
						TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
						ProjectID:      stringPtr("other_project_id"),
						TargetRepo:     stringPtr("other/repo"),
						Fields: map[string]string{
							"NonExistentField": "value",
						},
					},
				},
			},
			mockFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Status",
					DataType: "SINGLE_SELECT",
					Options: []ProjectFieldOption{
						{ID: "OPT_1", Name: "Ready"},
					},
				},
			},
			expectError:         true,
			expectErrorContains: "field 'NonExistentField' does not exist in project",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var mockClient *mockGitHubClient
			if tt.name == "invalid - field does not exist in overridden project" {
				// Use different fields for different projects
				mockClient = newMockGitHubClientWithMultipleProjects(map[string][]ProjectField{
					"default_project_id:default": {
						{
							ID:       "PVTFL_1",
							Name:     "Status",
							DataType: "SINGLE_SELECT",
							Options: []ProjectFieldOption{
								{ID: "OPT_1", Name: "Ready"},
							},
						},
					},
					"other_project_id:other": {
						{
							ID:       "PVTFL_2",
							Name:     "Priority",
							DataType: "TEXT",
						},
					},
				})
			} else {
				mockClient = newMockGitHubClient(tt.mockFields)
			}
			err := ValidateConfig(tt.config, mockClient)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectErrorContains)
					return
				}
				if tt.expectErrorContains != "" && !contains(err.Error(), tt.expectErrorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.expectErrorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// mockGitHubClient is a mock implementation of GitHubClient for testing
type mockGitHubClient struct {
	fieldsByProject map[string][]ProjectField
}

func (m *mockGitHubClient) GetProjectFields(ctx context.Context, projectID string, owner string) ([]ProjectField, error) {
	key := projectID + ":" + owner
	if fields, ok := m.fieldsByProject[key]; ok {
		return fields, nil
	}
	// Fallback to default project if key not found
	if fields, ok := m.fieldsByProject["default"]; ok {
		return fields, nil
	}
	return []ProjectField{}, nil
}

// newMockGitHubClient creates a mock GitHub client with fields for a single project
func newMockGitHubClient(fields []ProjectField) *mockGitHubClient {
	return &mockGitHubClient{
		fieldsByProject: map[string][]ProjectField{
			"default": fields,
		},
	}
}

// newMockGitHubClientWithMultipleProjects creates a mock GitHub client with fields for multiple projects
func newMockGitHubClientWithMultipleProjects(projectFields map[string][]ProjectField) *mockGitHubClient {
	return &mockGitHubClient{
		fieldsByProject: projectFields,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestValidateIssue(t *testing.T) {
	cases := []struct {
		name                string
		issue               Issue
		expectError         bool
		expectErrorContains string
	}{
		{
			name: "valid issue",
			issue: Issue{
				Name:           "test",
				CreationMonths: []Month{January},
				TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
			},
			expectError: false,
		},
		{
			name: "invalid - empty name",
			issue: Issue{
				Name:           "",
				CreationMonths: []Month{January},
				TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
			},
			expectError:         true,
			expectErrorContains: "name is required",
		},
		{
			name: "invalid - empty creation_months",
			issue: Issue{
				Name:           "test",
				CreationMonths: []Month{},
				TemplateFile:   stringPtr(".github/ISSUE_TEMPLATE/test.md"),
			},
			expectError:         true,
			expectErrorContains: "creation_months is required and must not be empty",
		},
		{
			name: "invalid - nil template_file",
			issue: Issue{
				Name:           "test",
				CreationMonths: []Month{January},
				TemplateFile:   nil,
			},
			expectError:         true,
			expectErrorContains: "template_file is required",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIssue(tt.issue)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectErrorContains)
					return
				}
				if tt.expectErrorContains != "" && !contains(err.Error(), tt.expectErrorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.expectErrorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
