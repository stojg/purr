package main

import "testing"

func TestFilters_Filter(t *testing.T) {

	tests := []struct {
		pr        *PullRequest
		whitelist []string
		expected  bool
	}{
		{pr: &PullRequest{Author: "john"}, whitelist: []string{}, expected: false},
		{pr: &PullRequest{Author: "john"}, whitelist: []string{"jane", "john"}, expected: false},
		{pr: &PullRequest{Author: "john"}, whitelist: []string{"jane"}, expected: true},
		{pr: &PullRequest{Assignee: "john"}, whitelist: []string{"jane", "john"}, expected: false},
		{pr: &PullRequest{Assignee: "john"}, whitelist: []string{"jane"}, expected: true},
	}
	for _, test := range tests {
		filters := &Filters{}
		filters.SetUsers(test.whitelist)
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("Expected isWhitelisted '%t', got '%t', whitelist %s, author: '%s' and assignee '%s'", test.expected, actual, test.whitelist, test.pr.Author, test.pr.Assignee)
		}
	}
}
