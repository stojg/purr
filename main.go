package main

import (
	"bytes"
	"fmt"
	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"io"
	"os"
	"strings"
)

type config struct {
	githubToken  string
	githubRepos  []string
	gitlabToken  string
	gitlabRepos  []string
	gitlabURL    string
	slackToken   string
	slackChannel string
}

func newConfig() *config {
	c := &config{}
	c.githubToken = os.Getenv("GITHUB_TOKEN")
	if os.Getenv("GITHUB_REPOS") != "" {
		c.githubRepos = strings.Split(os.Getenv("GITHUB_REPOS"), ",")
	}
	c.gitlabToken = os.Getenv("GITLAB_TOKEN")
	if os.Getenv("GITLAB_REPOS") != "" {
		c.gitlabRepos = strings.Split(os.Getenv("GITLAB_REPOS"), ",")
	}
	c.gitlabURL = os.Getenv("GITLAB_URL")
	c.slackToken = os.Getenv("SLACK_TOKEN")
	c.slackChannel = os.Getenv("SLACK_CHANNEL")
	return c
}

func (c *config) validate() []error {
	var errors []error
	if c.githubToken == "" {
		errors = append(errors, fmt.Errorf("err: no token defined in ENV 'GITHUB_TOKEN'"))
	}
	if c.slackToken == "" {
		errors = append(errors, fmt.Errorf("err: no token found in ENV 'SLACK_TOKEN'"))
	}
	if len(c.githubRepos) == 0 {
		errors = append(errors, fmt.Errorf("err: no repos defined in ENV 'GITHUB_REPOS'"))
	}
	if c.slackChannel == "" {
		errors = append(errors, fmt.Errorf("err: no slack channel found in ENV 'SLACK_CHANNEL'"))
	}

	return errors
}
func main() {
	conf := newConfig()

	if errors := conf.validate(); len(errors) > 0 {
		for i := range errors {
			fmt.Println(errors[i])
		}
		os.Exit(1)
	}

	var message []byte
	buf := bytes.NewBuffer(message)

	getGithub(conf, buf)
	getGitlab(conf, buf)
	sendToSlack(conf, buf)

}

func sendToSlack(conf *config, message fmt.Stringer) {

	if message.String() != "" {
		client := slack.New(conf.slackToken)
		opt := &slack.ChatPostMessageOpt{
			AsUser:    false,
			Username:  "purr",
			IconEmoji: ":purr:",
		}
		err := client.ChatPostMessage(conf.slackChannel, message.String(), opt)
		if err != nil {
			fmt.Printf("while sending slack request: %s\n", err)
			os.Exit(1)
		}
	}

}

func getGithub(conf *config, message io.Writer) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.githubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	for _, repo := range conf.githubRepos {
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

		if len(pullRequests) == 0 {
			continue
		}

		fmt.Fprintf(message, "*%s*\n", repo)
		for _, pr := range pullRequests {
			if isWIP(*pr.Title) {
				continue
			}
			timeAgo := humanize.Time(*pr.UpdatedAt)
			fmt.Fprintf(message, " • <%s|%s> - %s - updated %s\n", *pr.HTMLURL, *pr.Title, *pr.User.Login, timeAgo)
		}
		fmt.Fprintf(message, "\n")
	}
}

func getGitlab(conf *config, message io.Writer) {

	if conf.gitlabURL == "" {
		return
	}

	client := gitlab.NewClient(nil, conf.gitlabToken)
	client.SetBaseURL(conf.gitlabURL + "/api/v3")
	status := "opened"
	options := &gitlab.ListMergeRequestsOptions{
		State: &status,
	}

	for _, repo := range conf.gitlabRepos {
		pullRequests, _, err := client.MergeRequests.ListMergeRequests(repo, options)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		if len(pullRequests) == 0 {
			continue
		}

		fmt.Fprintf(message, "*%s*\n", repo)
		for _, pr := range pullRequests {
			if isWIP(pr.Title) {
				continue
			}
			timeAgo := humanize.Time(*pr.UpdatedAt)
			webLink := fmt.Sprintf("%s/%s/merge_requests/%d", conf.gitlabURL, repo, pr.IID)
			fmt.Fprintf(message, " • <%s|%s> - %s - updated %s\n", webLink, pr.Title, pr.Author.Username, timeAgo)
		}
		fmt.Fprintf(message, "\n")
	}
}

func isWIP(s string) bool {
	return strings.Contains(s, "[WIP]") || strings.Contains(s, "WIP:")
}
