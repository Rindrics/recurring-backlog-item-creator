package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v62/github"
)

type GitHubClient interface {
	GetProjectFields(ctx context.Context, projectID string, owner string) ([]ProjectField, error)
	GetProjectName(ctx context.Context, projectID string) (string, error)
}

type ProjectField struct {
	ID       string
	Name     string
	DataType string
	Options  []ProjectFieldOption
}

type ProjectFieldOption struct {
	ID   string
	Name string
}

type githubClient struct {
	client *github.Client
}

func NewGitHubClient() (GitHubClient, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	httpClient := &http.Client{
		Transport: &tokenTransport{
			token: token,
		},
	}

	client := github.NewClient(httpClient)
	return &githubClient{client: client}, nil
}

type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return http.DefaultTransport.RoundTrip(req)
}

func (g *githubClient) GetProjectFields(ctx context.Context, projectID string, owner string) ([]ProjectField, error) {
	var allFields []ProjectField
	cursor := ""
	hasNextPage := true

	for hasNextPage {
		query := fmt.Sprintf(`
			query($cursor: String) {
				node(id: "%s") {
					... on ProjectV2 {
						fields(first: 100, after: $cursor) {
							pageInfo {
								hasNextPage
								endCursor
							}
							nodes {
								... on ProjectV2Field {
									id
									name
									dataType
								}
								... on ProjectV2SingleSelectField {
									id
									name
									dataType
									options {
										id
										name
									}
								}
							}
						}
					}
				}
			}
		`, projectID)

		variables := map[string]interface{}{}
		if cursor != "" {
			variables["cursor"] = cursor
		}

		var result struct {
			Data struct {
				Node struct {
					Fields struct {
						PageInfo struct {
							HasNextPage bool   `json:"hasNextPage"`
							EndCursor   string `json:"endCursor"`
						} `json:"pageInfo"`
						Nodes []struct {
							ID       string `json:"id"`
							Name     string `json:"name"`
							DataType string `json:"dataType"`
							Options  []struct {
								ID   string `json:"id"`
								Name string `json:"name"`
							} `json:"options,omitempty"`
						} `json:"nodes"`
					} `json:"fields"`
				} `json:"node"`
			} `json:"data"`
			Errors []struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"errors,omitempty"`
		}

		req, err := g.client.NewRequest("POST", "/graphql", map[string]interface{}{
			"query":     query,
			"variables": variables,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create GraphQL request: %w", err)
		}

		resp, err := g.client.Do(ctx, req, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to execute GraphQL query: %w", err)
		}
		defer resp.Body.Close()

		Debugf("GraphQL response status: %d", resp.StatusCode)

		// Check for GraphQL errors
		if len(result.Errors) > 0 {
			errorMessages := make([]string, 0, len(result.Errors))
			for _, err := range result.Errors {
				errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", err.Type, err.Message))
				Debugf("GraphQL error: %s - %s", err.Type, err.Message)
			}
			// Provide more helpful error message for NOT_FOUND errors
			for _, err := range result.Errors {
				if err.Type == "NOT_FOUND" {
					return nil, fmt.Errorf("project not found (ID: %s). This may indicate: 1) The project ID is incorrect, 2) The token doesn't have access to this project, or 3) The project belongs to a different organization/user. GraphQL error: %s", projectID, err.Message)
				}
			}
			return nil, fmt.Errorf("GraphQL errors: %v", errorMessages)
		}

		Debugf("GraphQL response - hasNextPage: %v, endCursor: %s, nodes count: %d",
			result.Data.Node.Fields.PageInfo.HasNextPage,
			result.Data.Node.Fields.PageInfo.EndCursor,
			len(result.Data.Node.Fields.Nodes))

		if len(result.Data.Node.Fields.Nodes) == 0 && cursor == "" {
			// First page is empty - this might indicate a permissions issue
			Debugf("Warning: No fields found in first page. This might indicate a permissions issue or the project has no custom fields.")
			Debugf("Project ID: %s, Owner: %s", projectID, owner)
		}

		for _, node := range result.Data.Node.Fields.Nodes {
			field := ProjectField{
				ID:       node.ID,
				Name:     node.Name,
				DataType: node.DataType,
			}

			// Add options for single-select fields
			if len(node.Options) > 0 {
				field.Options = make([]ProjectFieldOption, 0, len(node.Options))
				for _, opt := range node.Options {
					field.Options = append(field.Options, ProjectFieldOption{
						ID:   opt.ID,
						Name: opt.Name,
					})
				}
			}

			allFields = append(allFields, field)
		}

		hasNextPage = result.Data.Node.Fields.PageInfo.HasNextPage
		cursor = result.Data.Node.Fields.PageInfo.EndCursor
	}

	return allFields, nil
}

func (g *githubClient) GetProjectName(ctx context.Context, projectID string) (string, error) {
	query := fmt.Sprintf(`
		query {
			node(id: "%s") {
				... on ProjectV2 {
					title
				}
			}
		}
	`, projectID)

	var result struct {
		Data struct {
			Node struct {
				Title string `json:"title"`
			} `json:"node"`
		} `json:"data"`
	}

	req, err := g.client.NewRequest("POST", "/graphql", map[string]interface{}{
		"query": query,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create GraphQL request: %w", err)
	}

	_, err = g.client.Do(ctx, req, &result)
	if err != nil {
		return "", fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	if result.Data.Node.Title == "" {
		return "", fmt.Errorf("project name is empty for project ID %s", projectID)
	}

	return result.Data.Node.Title, nil
}

func NewGitHubClientWithHTTPClient(httpClient *http.Client) GitHubClient {
	client := github.NewClient(httpClient)
	return &githubClient{client: client}
}
