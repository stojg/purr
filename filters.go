package main

import "strings"

type Filter interface {
	// Filter returns true if a PR should be kept and false if it should be discarded
	Filter(*PullRequest) bool
}

type Filters struct {
	Users           UserFilter   `json:"users"`
	WorkInProgress  WIPFilter    `json:"wip"`
	RequiresChanges ReviewFilter `json:"review"`
}

// Filter returns true if a PR should be kept and false if it should be discarded
func (f *Filters) Filter(p *PullRequest) bool {
	if !f.Users.Filter(p) {
		return false
	}
	if !f.WorkInProgress.Filter(p) {
		return false
	}
	if !f.RequiresChanges.Filter(p) {
		return false
	}
	return true
}

// UserFilter filters out any PRs that is not authored or assigned to a user
type UserFilter []string

// Filter returns true if a PR should be kept and false if it should be discarded
func (users UserFilter) Filter(p *PullRequest) bool {
	if len(users) == 0 {
		return true
	}
	for _, user := range users {
		if user == p.Author {
			return true
		}
		if user == p.Assignee {
			return true
		}
	}
	return false
}

// WIPFilter checks if the PR has been marked as Work In Progress, typically by prefixing the title with "WIP"
type WIPFilter bool

// Filter returns true if a PR should be kept and false if it should be discarded
func (enabled WIPFilter) Filter(p *PullRequest) bool {
	if !enabled {
		return true
	}

	if strings.Index(p.Title, "[WIP]") == 0 {
		return false
	}

	if strings.Index(p.Title, "WIP") == 0 {
		return false
	}
	return true
}

// ReviewFilter filters out any PR that had changes requested and haven't yet been approved
type ReviewFilter bool

// Filter returns true if a PR should be kept and false if it should be discarded
func (enabled ReviewFilter) Filter(p *PullRequest) bool {
	if !enabled {
		return true
	}

	return !p.RequiresChanges

}
