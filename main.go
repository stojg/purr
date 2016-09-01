package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "purr - %s\n"
	// VERSION is the binary version.
	VERSION = "v0.3.0"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "", "Read config from FILE")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
		configHelp()
	}
	flag.Parse()
}

func main() {
	conf := newConfig()

	if errors := conf.validate(); len(errors) > 0 {
		buf := &bytes.Buffer{}
		for i := range errors {
			fmt.Fprintln(buf, errors[i].Error())
		}
		usageAndExit(buf.String(), 1)
	}

	prs := make(map[string][]*PullRequest, 0)
	setGitHub(conf, prs)
	setGitLab(conf, prs)
	prs = filter(prs, conf)

	buf := &bytes.Buffer{}
	for repo, list := range prs {
		fmt.Fprintf(buf, "*%s*\n", repo)
		for i := range list {
			fmt.Fprintf(buf, "%s\n", list[i])
		}
		fmt.Fprint(buf, "\n")
	}
	sendToSlack(conf, buf)
}

func setGitHub(conf *Config, prs map[string][]*PullRequest) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GitHubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	for _, repo := range conf.GitHubRepos {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			fmt.Printf("%s is not a valid github repository, skipping\n", repo)
			continue
		}
		pullRequests, _, err := client.PullRequests.List(parts[0], parts[1], nil)
		if err != nil {
			fmt.Printf("while fetching pull requests from '%s': %s\n", repo, err)
			continue
		}
		var results []*PullRequest
		for _, pr := range pullRequests {
			pullRequest := &PullRequest{
				Author:     *pr.User.Login,
				TimeAgo:    humanize.Time(*pr.UpdatedAt),
				WebLink:    *pr.HTMLURL,
				Title:      *pr.Title,
				Repository: repo,
			}
			if pr.Assignee != nil {
				pullRequest.Assignee = *pr.Assignee.Login
			}
			results = append(results, pullRequest)
		}
		if len(results) > 0 {
			prs[repo] = results
		}
	}
}

func setGitLab(conf *Config, prs map[string][]*PullRequest) {
	client := gitlab.NewClient(nil, conf.GitLabToken)
	if err := client.SetBaseURL(conf.GitlabURL + "/api/v3"); err != nil {
		usageAndExit(err.Error(), 1)
	}
	status := "opened"
	options := &gitlab.ListMergeRequestsOptions{State: &status}
	for _, repo := range conf.GitLabRepos {
		pullRequests, _, err := client.MergeRequests.ListMergeRequests(repo, options)
		if err != nil {
			fmt.Printf("while fetching pull requests from '%s': %s\n", repo, err)
			continue
		}
		var results []*PullRequest
		for _, pr := range pullRequests {
			results = append(results, &PullRequest{
				Author:     pr.Author.Username,
				Assignee:   pr.Assignee.Username,
				TimeAgo:    humanize.Time(*pr.UpdatedAt),
				WebLink:    fmt.Sprintf("%s/%s/merge_requests/%d", conf.GitlabURL, repo, pr.IID),
				Title:      pr.Title,
				Repository: repo,
			})
		}
		if len(results) > 0 {
			prs[repo] = results
		}
	}
}

func filter(prs map[string][]*PullRequest, conf *Config) map[string][]*PullRequest {
	result := make(map[string][]*PullRequest)
	for repo, list := range prs {
		var filtered []*PullRequest
		for i := range list {
			if !isWIP(list[i].Title) && isWhiteListed(conf, list[i]) {
				filtered = append(filtered, list[i])
			}
		}
		if len(filtered) > 0 {
			result[repo] = filtered
		}
	}
	return result
}

func sendToSlack(conf *Config, message fmt.Stringer) {
	if message.String() != "" {
		client := slack.New(conf.SlackToken)
		opt := &slack.ChatPostMessageOpt{
			AsUser:    false,
			Username:  "purr",
			IconEmoji: ":purr:",
		}
		err := client.ChatPostMessage(conf.SlackChannel, message.String(), opt)
		if err != nil {
			fmt.Printf("while sending slack request: %s\n", err)
			os.Exit(1)
		}
	}
}

// PullRequest is a normalised version of PullRequest from the different providers
type PullRequest struct {
	Author     string
	Assignee   string
	TimeAgo    string
	WebLink    string
	Title      string
	Repository string
}

func (p *PullRequest) String() string {
	output := fmt.Sprintf(" • <%s|%s> - by %s", p.WebLink, p.Title, p.Author)
	if p.Assignee != "" {
		output += fmt.Sprintf(", assigned to %s", p.Assignee)
	}
	output += fmt.Sprintf(" - updated %s", p.TimeAgo)
	return output
}

func isWIP(s string) bool {
	return strings.Contains(s, "[WIP]") || strings.Contains(s, "WIP:")
}

func isWhiteListed(config *Config, pr *PullRequest) bool {
	if len(config.UserWhiteList) == 0 {
		return true
	}
	for _, user := range config.UserWhiteList {
		if user == pr.Author || user == pr.Assignee {
			return true
		}
	}
	return false
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
