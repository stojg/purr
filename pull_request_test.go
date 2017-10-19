package main

import (
	"testing"
)

func TestPullRequest_isWIP(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected bool
	}{
		{pr: &PullRequest{Title: "fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "something WIP fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "something WIP: fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "something [WIP] fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "SWIP fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "wip fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "[wip] fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "WIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "WIP: fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "[WIP] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "Title", HasChangesRequestedReview: true}, expected: true},
	}
	for _, test := range tests {
		if test.pr.isWIP() != test.expected {
			t.Errorf("Expected '%t', got '%t' for title %s", test.expected, test.pr.isWIP(), test.pr.Title)
		}
	}
}

func TestPullRequest_String(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected string
	}{
		{
			pr: &PullRequest{
				WebLink: "http://gitlab.local/324",
				ID:      324,
				Title:   "fixes bug",
				Author:  "john.doe",
			},
			expected: " • <http://gitlab.local/324|PR #324> fixes bug - _john.doe_ - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink: "http://gitlab.local/324",
				ID:      324,
				Title:   "fixes < bug with > & other",
				Author:  "john.doe",
			},
			expected: " • <http://gitlab.local/324|PR #324> fixes &lt; bug with &gt; &amp; other - _john.doe_ - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink:           "http://gitlab.local/324",
				ID:                324,
				Title:             "fixes bug",
				Author:            "john.doe",
				HasApprovedReview: true,
			},
			expected: " • <http://gitlab.local/324|PR #324> fixes bug - _john.doe_, *APPROVED* - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink:  "http://gitlab.local/243",
				ID:       243,
				Title:    "fixes bug",
				Author:   "john.doe",
				Assignee: "jane.doe",
			},
			expected: " • <http://gitlab.local/243|PR #243> fixes bug - _john.doe_, assigned to _jane.doe_ - updated a long while ago",
		},
	}

	for _, test := range tests {
		if test.pr.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.pr.String())
		}
	}
}
