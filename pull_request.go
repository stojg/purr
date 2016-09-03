package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
	"time"
)

// PullRequest is a normalised version of PullRequest for the different providers
type PullRequest struct {
	ID         int
	Author     string
	Assignee   string
	Updated    time.Time
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

	output := fmt.Sprintf(" • <%s|PR #%d> %s  - _%s_", p.WebLink, p.ID, p.Title, p.Author)
	if p.Assignee != "" {
		output += fmt.Sprintf(", assigned to %s", p.Assignee)
	}
	output += fmt.Sprintf(" - updated %s", humanize.Time(p.Updated))
	return output
}
