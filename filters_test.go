package main

import "testing"

func TestFilters_Filter(t *testing.T) {

	tests := []struct {
		pr       *PullRequest
		users    []string
		expected bool
	}{
		{pr: &PullRequest{Author: "john"}, users: []string{}, expected: false},
		{pr: &PullRequest{Author: "john"}, users: []string{"jane", "john"}, expected: false},
		{pr: &PullRequest{Author: "john"}, users: []string{"jane"}, expected: true},
		{pr: &PullRequest{Assignee: "john"}, users: []string{"jane", "john"}, expected: false},
		{pr: &PullRequest{Assignee: "john"}, users: []string{"jane"}, expected: true},
	}
	for _, test := range tests {
		filters := &Filters{}
		filters.SetUsers(test.users)
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("Expected Filter '%t', got '%t', users %s", test.expected, actual, test.users)
		}
	}
}
