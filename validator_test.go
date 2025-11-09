package main

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	cases := []struct {
		name                string
		config              Config
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
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid - empty project_id",
			config: Config{
				Defaults: Defaults{
					ProjectID: "",
				},
			},
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
			expectError:         true,
			expectErrorContains: "name is required",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
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
