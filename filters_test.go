package main

import "testing"

type mockFilter bool

func (m mockFilter) Filter(p *PullRequest) bool {
	return bool(m)
}

func TestFilter_NumFilteredNone(t *testing.T) {
	filters := &Filters{}
	filters.Add(mockFilter(true))

	filters.Filter(&PullRequest{})
	filters.Filter(&PullRequest{})
	filters.Filter(&PullRequest{})

	if filters.NumFiltered() != 0 {
		t.Errorf("Expected 0 filtered events, got %d", filters.NumFiltered())
	}
}

func TestFilter_NumFilteredAll(t *testing.T) {
	filters := &Filters{}
	filters.Add(mockFilter(false))

	filters.Filter(&PullRequest{})
	filters.Filter(&PullRequest{})
	filters.Filter(&PullRequest{})

	if filters.NumFiltered() != 3 {
		t.Errorf("Expected 3 filtered events, got %d", filters.NumFiltered())
	}
}

func TestFilters_FilterUsers(t *testing.T) {

	tests := []struct {
		pr       *PullRequest
		users    []string
		expected bool
	}{
		{pr: &PullRequest{Author: "john"}, users: []string{}, expected: true},
		{pr: &PullRequest{Author: "john"}, users: []string{"jane", "john"}, expected: true},
		{pr: &PullRequest{Author: "john"}, users: []string{"jane"}, expected: false},
		{pr: &PullRequest{Assignee: "john"}, users: []string{"jane", "john"}, expected: true},
		{pr: &PullRequest{Assignee: "john"}, users: []string{"jane"}, expected: false},
	}
	for i, test := range tests {
		filters := &Filters{}
		filters.Add(UserFilter(test.users))
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("case %d. Expected Filter '%t', got '%t', users %s, author: '%s', assigned: '%s'", i+1, test.expected, actual, test.users, test.pr.Author, test.pr.Assignee)
		}
	}
}

func TestWorkInProgressFilter_Filter(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected bool
	}{
		{pr: &PullRequest{Title: "fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something WIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something WIP: fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something [WIP] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "SWIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "wip fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "[wip] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "WIP fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "WIP: fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "[WIP] fixes bug"}, expected: false},
		{pr: &PullRequest{Title: "fixes bug", Draft: true}, expected: false},
		{pr: &PullRequest{Title: "fixes bug", Draft: false}, expected: true},
	}

	for _, test := range tests {
		filters := &Filters{}
		filters.Add(WIPFilter(true))
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("Expected '%t', got '%t' for title %s", test.expected, actual, test.pr.Title)
		}
	}
}

func TestWorkInProgressFilter_FilterDisabled(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected bool
	}{
		{pr: &PullRequest{Title: "fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something WIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something WIP: fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "something [WIP] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "SWIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "wip fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "[wip] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "WIP fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "WIP: fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "[WIP] fixes bug"}, expected: true},
		{pr: &PullRequest{Title: "fixes bug", Draft: true}, expected: true},
		{pr: &PullRequest{Title: "fixes bug", Draft: false}, expected: true},
	}

	for _, test := range tests {
		filters := &Filters{}
		filters.Add(WIPFilter(false))
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("Expected '%t', got '%t' for title %s", test.expected, actual, test.pr.Title)
		}
	}
}

func TestReviewFilter_Filter(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected bool
	}{
		{pr: &PullRequest{}, expected: true},
		{pr: &PullRequest{RequiresChanges: true}, expected: false},
		{pr: &PullRequest{RequiresChanges: false}, expected: true},
	}

	for _, test := range tests {
		filters := &Filters{}
		filters.Add(ReviewFilter(true))
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("Expected '%t', got '%t'", test.expected, actual)
		}
	}
}

func TestReviewFilter_FilterDisabled(t *testing.T) {
	tests := []struct {
		pr       *PullRequest
		expected bool
	}{
		{pr: &PullRequest{RequiresChanges: false}, expected: true},
		{pr: &PullRequest{RequiresChanges: true}, expected: true},
	}

	for i, test := range tests {
		filters := &Filters{}
		filters.Add(ReviewFilter(false))
		actual := filters.Filter(test.pr)
		if actual != test.expected {
			t.Errorf("case %d. Expected '%t', got '%t'", i+1, test.expected, actual)
		}
	}
}
