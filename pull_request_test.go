package main

import (
	"testing"
)

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
			expected: " • <http://gitlab.local/324|#324> fixes bug - _john.doe_ - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink: "http://gitlab.local/324",
				ID:      324,
				Title:   "fixes < bug with > & other",
				Author:  "john.doe",
			},
			expected: " • <http://gitlab.local/324|#324> fixes &lt; bug with &gt; &amp; other - _john.doe_ - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink:  "http://gitlab.local/324",
				ID:       324,
				Title:    "fixes bug",
				Author:   "john.doe",
				Approved: true,
			},
			expected: " • <http://gitlab.local/324|#324> fixes bug - _john.doe_, *APPROVED* - updated a long while ago",
		},
		{
			pr: &PullRequest{
				WebLink:  "http://gitlab.local/243",
				ID:       243,
				Title:    "fixes bug",
				Author:   "john.doe",
				Assignee: "jane.doe",
			},
			expected: " • <http://gitlab.local/243|#243> fixes bug - _john.doe_, assigned to _jane.doe_ - updated a long while ago",
		},
	}

	for _, test := range tests {
		if test.pr.String() != test.expected {
			t.Errorf("Expected '%s', got '%s'", test.expected, test.pr.String())
		}
	}
}
