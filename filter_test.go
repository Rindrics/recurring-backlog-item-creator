package main

import (
	"testing"
)

func TestGetIssuesToCreate(t *testing.T) {
	defaults := Defaults{
		ProjectID:  "default_project_id",
		TargetRepo: "default/repo",
	}

	stringPtr := func(s string) *string {
		return &s
	}
	issue1 := Issue{
		Name:           "Issue 1",
		CreationMonths: []Month{January},
	}
	issue2 := Issue{
		Name:           "Issue 2",
		CreationMonths: []Month{February},
	}
	issue1_3 := Issue{
		Name:           "Issue 1_3",
		CreationMonths: []Month{January, March},
	}
	issue2_4 := Issue{
		Name:           "Issue 2_4",
		CreationMonths: []Month{February, April},
	}
	otherProjectID := "other_project_id"
	otherRepo := "other/repo"
	issue_project_repo := Issue{
		Name:           "Issue project_repo",
		CreationMonths: []Month{January},
		ProjectID:      &otherProjectID,
		TargetRepo:     &otherRepo,
	}

	cases := []struct {
		name           string
		config         Config
		month          Month
		issuesToCreate IssuesToCreate
	}{
		{
			name: "No issues",
			config: Config{
				Defaults: defaults,
				Issues:   []Issue{},
			},
			month:          January,
			issuesToCreate: IssuesToCreate{},
		},
		{
			name: "One issue",
			config: Config{
				Defaults: defaults,
				Issues:   []Issue{issue1},
			},
			month: January,
			issuesToCreate: IssuesToCreate{
				Issues: []IssueToCreate{
					{
						Issue:      issue1,
						ProjectID:  stringPtr("default_project_id"),
						TargetRepo: stringPtr("default/repo"),
					},
				},
			},
		},
		{
			name: "January issues",
			config: Config{
				Defaults: defaults,
				Issues:   []Issue{issue1, issue1_3, issue2_4},
			},
			month: January,
			issuesToCreate: IssuesToCreate{
				Issues: []IssueToCreate{
					{
						Issue:      issue1,
						ProjectID:  stringPtr("default_project_id"),
						TargetRepo: stringPtr("default/repo"),
					},
					{
						Issue:      issue1_3,
						ProjectID:  stringPtr("default_project_id"),
						TargetRepo: stringPtr("default/repo"),
					},
				},
			},
		},
		{
			name: "February issues",
			config: Config{
				Defaults: defaults,
				Issues:   []Issue{issue2, issue1_3, issue2_4},
			},
			month: February,
			issuesToCreate: IssuesToCreate{
				Issues: []IssueToCreate{
					{
						Issue:      issue2,
						ProjectID:  stringPtr("default_project_id"),
						TargetRepo: stringPtr("default/repo"),
					},
					{
						Issue:      issue2_4,
						ProjectID:  stringPtr("default_project_id"),
						TargetRepo: stringPtr("default/repo"),
					},
				},
			},
		},
		{
			name: "Override defaults",
			config: Config{
				Defaults: defaults,
				Issues:   []Issue{issue_project_repo},
			},
			month: January,
			issuesToCreate: IssuesToCreate{
				Issues: []IssueToCreate{
					{
						Issue:      issue_project_repo,
						ProjectID:  &otherProjectID,
						TargetRepo: &otherRepo,
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIssuesToCreate(tt.config, tt.month)
			if !got.Equals(tt.issuesToCreate) {
				t.Errorf("expected %v, got %v", tt.issuesToCreate, got)
			}
		})
	}
}
