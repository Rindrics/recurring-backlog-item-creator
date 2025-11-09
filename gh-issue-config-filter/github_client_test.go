package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjectFields(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		owner          string
		mockResponse   map[string]interface{}
		expectedFields []ProjectField
		expectError    bool
	}{
		{
			name:      "successfully get project fields",
			projectID: "PVT_kwHOAOKHl84BHgin",
			owner:     "test-owner",
			mockResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"node": map[string]interface{}{
						"fields": map[string]interface{}{
							"pageInfo": map[string]interface{}{
								"hasNextPage": false,
								"endCursor":   nil,
							},
							"nodes": []interface{}{
								map[string]interface{}{
									"id":       "PVTFL_1",
									"name":     "Status",
									"dataType": "SINGLE_SELECT",
									"options": []interface{}{
										map[string]interface{}{
											"id":   "OPT_1",
											"name": "Ready",
										},
										map[string]interface{}{
											"id":   "OPT_2",
											"name": "In Progress",
										},
									},
								},
								map[string]interface{}{
									"id":       "PVTFL_2",
									"name":     "Story Points",
									"dataType": "NUMBER",
								},
							},
						},
					},
				},
			},
			expectedFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Status",
					DataType: "SINGLE_SELECT",
					Options: []ProjectFieldOption{
						{ID: "OPT_1", Name: "Ready"},
						{ID: "OPT_2", Name: "In Progress"},
					},
				},
				{
					ID:       "PVTFL_2",
					Name:     "Story Points",
					DataType: "NUMBER",
					Options:  []ProjectFieldOption{},
				},
			},
			expectError: false,
		},
		{
			name:      "handle pagination",
			projectID: "PVT_kwHOAOKHl84BHgin",
			owner:     "test-owner",
			mockResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"node": map[string]interface{}{
						"fields": map[string]interface{}{
							"pageInfo": map[string]interface{}{
								"hasNextPage": true,
								"endCursor":   "cursor1",
							},
							"nodes": []interface{}{
								map[string]interface{}{
									"id":       "PVTFL_1",
									"name":     "Field1",
									"dataType": "TEXT",
								},
							},
						},
					},
				},
			},
			expectedFields: []ProjectField{
				{
					ID:       "PVTFL_1",
					Name:     "Field1",
					DataType: "TEXT",
					Options:  []ProjectFieldOption{},
				},
				{
					ID:       "PVTFL_2",
					Name:     "Field2",
					DataType: "TEXT",
					Options:  []ProjectFieldOption{},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				if r.Method != "POST" {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/graphql" {
					t.Errorf("expected /graphql path, got %s", r.URL.Path)
				}

				// Parse request body
				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				// Handle pagination: if cursor is provided, return second page
				variables, ok := reqBody["variables"].(map[string]interface{})
				if ok && variables["cursor"] != nil {
					// Second page response
					response := map[string]interface{}{
						"data": map[string]interface{}{
							"node": map[string]interface{}{
								"fields": map[string]interface{}{
									"pageInfo": map[string]interface{}{
										"hasNextPage": false,
										"endCursor":   nil,
									},
									"nodes": []interface{}{
										map[string]interface{}{
											"id":       "PVTFL_2",
											"name":     "Field2",
											"dataType": "TEXT",
										},
									},
								},
							},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
					return
				}

				// First page response
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create HTTP client that points to mock server
			httpClient := &http.Client{
				Transport: &mockTransport{
					baseURL: server.URL,
				},
			}

			// Create GitHub client with mock HTTP client
			client := NewGitHubClientWithHTTPClient(httpClient)

			// Call GetProjectFields
			fields, err := client.GetProjectFields(context.Background(), tt.projectID, tt.owner)

			// Check error
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check fields
			if len(fields) != len(tt.expectedFields) {
				t.Errorf("expected %d fields, got %d", len(tt.expectedFields), len(fields))
			}

			for i, expected := range tt.expectedFields {
				if i >= len(fields) {
					break
				}
				actual := fields[i]
				if actual.ID != expected.ID {
					t.Errorf("field[%d].ID: expected %s, got %s", i, expected.ID, actual.ID)
				}
				if actual.Name != expected.Name {
					t.Errorf("field[%d].Name: expected %s, got %s", i, expected.Name, actual.Name)
				}
				if actual.DataType != expected.DataType {
					t.Errorf("field[%d].DataType: expected %s, got %s", i, expected.DataType, actual.DataType)
				}
				if len(actual.Options) != len(expected.Options) {
					t.Errorf("field[%d].Options: expected %d options, got %d", i, len(expected.Options), len(actual.Options))
				}
				for j, expectedOpt := range expected.Options {
					if j >= len(actual.Options) {
						break
					}
					actualOpt := actual.Options[j]
					if actualOpt.ID != expectedOpt.ID {
						t.Errorf("field[%d].Options[%d].ID: expected %s, got %s", i, j, expectedOpt.ID, actualOpt.ID)
					}
					if actualOpt.Name != expectedOpt.Name {
						t.Errorf("field[%d].Options[%d].Name: expected %s, got %s", i, j, expectedOpt.Name, actualOpt.Name)
					}
				}
			}
		})
	}
}

// mockTransport redirects requests to the mock server
type mockTransport struct {
	baseURL string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite URL to point to mock server
	req.URL.Scheme = "http"
	// Extract host and port from baseURL (e.g., "http://127.0.0.1:12345" -> "127.0.0.1:12345")
	host := m.baseURL[len("http://"):]
	req.URL.Host = host
	return http.DefaultTransport.RoundTrip(req)
}
