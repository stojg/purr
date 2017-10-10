package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// PullRequest is a normalised version of PullRequest for the different providers
type PullRequest struct {
	ID                        int
	Author                    string
	Assignee                  string
	Updated                   time.Time
	WebLink                   string
	Title                     string
	Repository                string
	HasChangesRequestedReview bool
	HasApprovedReview         bool
}

func (p *PullRequest) isWIP() bool {
	// Check for WIP in the title
	if strings.Contains(p.Title, "[WIP]") || strings.Contains(p.Title, "WIP:") {
		return true
	}

	// Check for a Review marked as requested changes
	if p.HasChangesRequestedReview {
		return true
	}

	return false
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

	// the following chars will to be escaped for slack
	title := strings.Replace(p.Title, "&", "&amp;", -1)
	title = strings.Replace(title, "<", "&lt;", -1)
	title = strings.Replace(title, ">", "&gt;", -1)

	output := fmt.Sprintf(" • <%s|PR #%d> %s  - _%s_", p.WebLink, p.ID, title, p.Author)

	if p.HasApprovedReview {
		output += ", *APPROVED*"
	}

	if p.Assignee != "" {
		output += fmt.Sprintf(", assigned to _%s_", p.Assignee)
	}
	output += fmt.Sprintf(" - updated %s", humanize.Time(p.Updated))
	return output
}
