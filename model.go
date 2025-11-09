package main

import (
	"reflect"
	"slices"
)

type Defaults struct {
	ProjectID  string `yaml:"project_id"`
	TargetRepo string `yaml:"target_repo"`
}

type Config struct {
	Defaults Defaults `yaml:"defaults"`
	Issues   []Issue  `yaml:"issues"`
}

type Month int

const (
	January   Month = 1
	February  Month = 2
	March     Month = 3
	April     Month = 4
	May       Month = 5
	June      Month = 6
	July      Month = 7
	August    Month = 8
	September Month = 9
	October   Month = 10
	November  Month = 11
	December  Month = 12
)

type Issue struct {
	Name           string            `yaml:"name"`
	CreationMonths []Month           `yaml:"creation_months"`
	TemplateFile   *string           `yaml:"template_file"`
	TitleSuffix    *string           `yaml:"title_suffix,omitempty"`
	Fields         map[string]string `yaml:"fields"`
	ProjectID      *string           `yaml:"project_id,omitempty"`
	TargetRepo     *string           `yaml:"target_repo,omitempty"`
}

type IssueToCreate = Issue

func NewIssueToCreate(issue Issue, defaults Defaults) IssueToCreate {
	issueToCreate := issue

	if issue.ProjectID == nil {
		projectID := defaults.ProjectID
		issueToCreate.ProjectID = &projectID
	}

	if issue.TargetRepo == nil {
		targetRepo := defaults.TargetRepo
		issueToCreate.TargetRepo = &targetRepo
	}

	return issueToCreate
}

type IssuesToCreate struct {
	Issues []IssueToCreate
}

func (i *IssuesToCreate) Equals(other IssuesToCreate) bool {
	if len(i.Issues) == 0 && len(other.Issues) == 0 {
		return true
	}
	return reflect.DeepEqual(i.Issues, other.Issues)
}

func (i *Issue) IsCreationMonth(month Month) bool {
	return slices.Contains(i.CreationMonths, month)
}
