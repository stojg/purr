package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// PullRequest is a normalised version of PullRequest for the different providers
type PullRequest struct {
	ID              int
	Author          string
	Assignee        string
	Updated         time.Time
	WebLink         string
	Title           string
	Repository      string
	RequiresChanges bool
	Approved        bool
	Draft           bool
}

func (p *PullRequest) String() string {

	// the following chars will to be escaped for slack
	title := strings.Replace(p.Title, "&", "&amp;", -1)
	title = strings.Replace(title, "<", "&lt;", -1)
	title = strings.Replace(title, ">", "&gt;", -1)

	output := fmt.Sprintf(" • <%s|#%d> %s - _%s_", p.WebLink, p.ID, title, p.Author)

	if p.Approved {
		output += ", *APPROVED*"
	}

	if p.Assignee != "" {
		output += fmt.Sprintf(", assigned to _%s_", p.Assignee)
	}
	output += fmt.Sprintf(" - updated %s", humanize.Time(p.Updated))
	return output
}
