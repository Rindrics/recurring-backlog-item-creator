package main

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestExpandTitlePrefix(t *testing.T) {
	now := time.Now()
	expectedYear := now.Format("2006")
	expectedMonth := now.Format("01")
	expectedYearMonth := now.Format("2006-01")
	expectedDate := now.Format("2006-01-02")

	cases := []struct {
		name        string
		titlePrefix *string
		expect      string
		expectError bool
	}{
		{
			name:        "nil prefix",
			titlePrefix: nil,
			expect:      "",
			expectError: false,
		},
		{
			name:        "empty prefix",
			titlePrefix: stringPtr(""),
			expect:      "",
			expectError: false,
		},
		{
			name:        "simple prefix",
			titlePrefix: stringPtr("[PREFIX]"),
			expect:      "[PREFIX]",
			expectError: false,
		},
		{
			name:        "prefix with Year",
			titlePrefix: stringPtr("{{Year}}"),
			expect:      expectedYear,
			expectError: false,
		},
		{
			name:        "prefix with Month",
			titlePrefix: stringPtr("{{Month}}"),
			expect:      expectedMonth,
			expectError: false,
		},
		{
			name:        "prefix with YearMonth",
			titlePrefix: stringPtr("{{YearMonth}}"),
			expect:      expectedYearMonth,
			expectError: false,
		},
		{
			name:        "prefix with Date",
			titlePrefix: stringPtr("{{Date}}"),
			expect:      expectedDate,
			expectError: false,
		},
		{
			name:        "prefix with multiple template variables",
			titlePrefix: stringPtr("[{{YearMonth}}]"),
			expect:      "[" + expectedYearMonth + "]",
			expectError: false,
		},
		{
			name:        "invalid template",
			titlePrefix: stringPtr("{{Invalid}}"),
			expect:      "",
			expectError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTitlePrefix(tt.titlePrefix)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.expect {
					t.Errorf("expected %q, got %q", tt.expect, got)
				}
			}
		})
	}
}

func TestExpandTitleSuffix(t *testing.T) {
	now := time.Now()
	expectedYear := now.Format("2006")
	expectedMonth := now.Format("01")
	expectedYearMonth := now.Format("2006-01")
	expectedDate := now.Format("2006-01-02")

	cases := []struct {
		name        string
		titleSuffix *string
		expect      string
		expectError bool
	}{
		{
			name:        "nil suffix",
			titleSuffix: nil,
			expect:      "",
			expectError: false,
		},
		{
			name:        "empty suffix",
			titleSuffix: stringPtr(""),
			expect:      "",
			expectError: false,
		},
		{
			name:        "simple suffix",
			titleSuffix: stringPtr("- suffix"),
			expect:      "- suffix",
			expectError: false,
		},
		{
			name:        "suffix with Year",
			titleSuffix: stringPtr("- {{Year}}"),
			expect:      "- " + expectedYear,
			expectError: false,
		},
		{
			name:        "suffix with Month",
			titleSuffix: stringPtr("- {{Month}}"),
			expect:      "- " + expectedMonth,
			expectError: false,
		},
		{
			name:        "suffix with YearMonth",
			titleSuffix: stringPtr("- {{YearMonth}}"),
			expect:      "- " + expectedYearMonth,
			expectError: false,
		},
		{
			name:        "suffix with Date",
			titleSuffix: stringPtr("- {{Date}}"),
			expect:      "- " + expectedDate,
			expectError: false,
		},
		{
			name:        "suffix with multiple template variables",
			titleSuffix: stringPtr("- ({{Year}})"),
			expect:      "- (" + expectedYear + ")",
			expectError: false,
		},
		{
			name:        "invalid template",
			titleSuffix: stringPtr("{{Invalid}}"),
			expect:      "",
			expectError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTitleSuffix(tt.titleSuffix)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.expect {
					t.Errorf("expected %q, got %q", tt.expect, got)
				}
			}
		})
	}
}

func TestBuildTitleWithPrefixAndSuffix(t *testing.T) {
	now := time.Now()
	expectedYear := now.Format("2006")
	expectedYearMonth := now.Format("2006-01")

	cases := []struct {
		name        string
		issueName   string
		titlePrefix *string
		titleSuffix *string
		expect      string
		expectError bool
	}{
		{
			name:        "no prefix, no suffix",
			issueName:   "Test Issue",
			titlePrefix: nil,
			titleSuffix: nil,
			expect:      "Test Issue",
			expectError: false,
		},
		{
			name:        "prefix only",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[PREFIX]"),
			titleSuffix: nil,
			expect:      "[PREFIX] Test Issue",
			expectError: false,
		},
		{
			name:        "suffix only",
			issueName:   "Test Issue",
			titlePrefix: nil,
			titleSuffix: stringPtr("- suffix"),
			expect:      "Test Issue - suffix",
			expectError: false,
		},
		{
			name:        "both prefix and suffix",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[PREFIX]"),
			titleSuffix: stringPtr("- suffix"),
			expect:      "[PREFIX] Test Issue - suffix",
			expectError: false,
		},
		{
			name:        "prefix and suffix with template variables",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[{{YearMonth}}]"),
			titleSuffix: stringPtr("- {{Year}}"),
			expect:      "[" + expectedYearMonth + "] Test Issue - " + expectedYear,
			expectError: false,
		},
		{
			name:        "empty prefix, suffix with template",
			issueName:   "Test Issue",
			titlePrefix: stringPtr(""),
			titleSuffix: stringPtr("- {{YearMonth}}"),
			expect:      "Test Issue - " + expectedYearMonth,
			expectError: false,
		},
		{
			name:        "prefix with template, empty suffix",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[{{YearMonth}}]"),
			titleSuffix: stringPtr(""),
			expect:      "[" + expectedYearMonth + "] Test Issue",
			expectError: false,
		},
		{
			name:        "prefix ending with space",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[PREFIX] "),
			titleSuffix: nil,
			expect:      "[PREFIX] Test Issue",
			expectError: false,
		},
		{
			name:        "suffix starting with space",
			issueName:   "Test Issue",
			titlePrefix: nil,
			titleSuffix: stringPtr(" - suffix"),
			expect:      "Test Issue - suffix",
			expectError: false,
		},
		{
			name:        "prefix ending with space and suffix starting with space",
			issueName:   "Test Issue",
			titlePrefix: stringPtr("[PREFIX] "),
			titleSuffix: stringPtr(" - suffix"),
			expect:      "[PREFIX] Test Issue - suffix",
			expectError: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			expandedPrefix, err := expandTitlePrefix(tt.titlePrefix)
			if err != nil && !tt.expectError {
				t.Fatalf("unexpected error expanding prefix: %v", err)
			}
			if err == nil && tt.expectError {
				// If we expect an error, it should come from suffix expansion
			}

			expandedSuffix, err := expandTitleSuffix(tt.titleSuffix)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error expanding suffix: %v", err)
			}

			// Build title: "{prefix} {name} {suffix}"
			// Add space between prefix and name only if prefix doesn't end with space
			// Add space between name and suffix only if suffix doesn't start with space
			title := tt.issueName
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

			if title != tt.expect {
				t.Errorf("expected %q, got %q", tt.expect, title)
			}
		})
	}
}

func TestExpandTitleTemplate(t *testing.T) {
	now := time.Now()
	expectedYear := now.Format("2006")

	cases := []struct {
		name         string
		templateStr  *string
		templateName string
		expect       string
		expectError  bool
		errorPattern *regexp.Regexp
	}{
		{
			name:         "nil template",
			templateStr:  nil,
			templateName: "test",
			expect:       "",
			expectError:  false,
		},
		{
			name:         "empty template",
			templateStr:  stringPtr(""),
			templateName: "test",
			expect:       "",
			expectError:  false,
		},
		{
			name:         "simple template",
			templateStr:  stringPtr("simple"),
			templateName: "test",
			expect:       "simple",
			expectError:  false,
		},
		{
			name:         "template with Year",
			templateStr:  stringPtr("{{Year}}"),
			templateName: "test",
			expect:       expectedYear,
			expectError:  false,
		},
		{
			name:         "invalid template function",
			templateStr:  stringPtr("{{Invalid}}"),
			templateName: "test",
			expect:       "",
			expectError:  true,
			errorPattern: regexp.MustCompile("failed to parse test template"),
		},
		{
			name:         "invalid template syntax",
			templateStr:  stringPtr("{{"),
			templateName: "test",
			expect:       "",
			expectError:  true,
			errorPattern: regexp.MustCompile("failed to parse test template"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTitleTemplate(tt.templateStr, tt.templateName)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errorPattern != nil && !tt.errorPattern.MatchString(err.Error()) {
					t.Errorf("error message %q does not match pattern %q", err.Error(), tt.errorPattern.String())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.expect {
					t.Errorf("expected %q, got %q", tt.expect, got)
				}
			}
		})
	}
}
