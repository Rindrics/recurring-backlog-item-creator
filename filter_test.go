package main

import (
	"testing"
)

func TestGetIssuesToCreate(t *testing.T) {
	defaults := Defaults{
		ProjectID:  "default_project_id",
		TargetRepo: "default/repo",
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
					NewIssueToCreate(issue1, defaults),
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
					NewIssueToCreate(issue1, defaults),
					NewIssueToCreate(issue1_3, defaults),
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
					NewIssueToCreate(issue2, defaults),
					NewIssueToCreate(issue2_4, defaults),
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
					NewIssueToCreate(issue_project_repo, defaults),
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
