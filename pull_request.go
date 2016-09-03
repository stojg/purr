package main

import (
	"fmt"
	"strings"
)

// PullRequest is a normalised version of PullRequest for the different providers
type PullRequest struct {
	Author     string
	Assignee   string
	TimeAgo    string
	WebLink    string
	Title      string
	Repository string
}

func (p *PullRequest) isWIP() bool {
	return strings.Contains(p.Title, "[WIP]") || strings.Contains(p.Title, "WIP:")
}

func (p *PullRequest) isWhiteListed(config *Config) bool {
	if len(config.UserWhiteList) == 0 {
		return true
	}
	for _, user := range config.UserWhiteList {
		if user == p.Author || user == p.Assignee {
			return true
		}
	}
	return false
}

func (p *PullRequest) String() string {
	output := fmt.Sprintf(" • <%s|%s> - by %s", p.WebLink, p.Title, p.Author)
	if p.Assignee != "" {
		output += fmt.Sprintf(", assigned to %s", p.Assignee)
	}
	output += fmt.Sprintf(" - updated %s", p.TimeAgo)
	return output
}
