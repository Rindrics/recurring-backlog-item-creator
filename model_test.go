package main

import (
	"testing"
)

func TestIssue_IsCreationMonth(t *testing.T) {

	cases := []struct {
		name           string
		creationMonths []int
		thisMonth      int
		expect         bool
	}{
		{
			name:           "Matches first month",
			creationMonths: []int{1, 2, 3},
			thisMonth:      1,
			expect:         true,
		},
		{
			name:           "Matches second month",
			creationMonths: []int{1, 2, 3},
			thisMonth:      2,
			expect:         true,
		},
		{
			name:           "Does not match",
			creationMonths: []int{1, 2, 3},
			thisMonth:      4,
			expect:         false,
		},
		{
			name:           "Empty creation months",
			creationMonths: []int{},
			thisMonth:      3,
			expect:         false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			issue := Issue{
				CreationMonths: tt.creationMonths,
			}
			if issue.IsCreationMonth(tt.thisMonth) != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, issue.IsCreationMonth(tt.thisMonth))
			}
		})
	}
}
