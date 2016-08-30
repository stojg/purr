package main

import (
	"fmt"
	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

type config struct {
	githubToken  string
	githubRepos  []string
	slackToken   string
	slackChannel string
}

func newConfig() *config {
	c := &config{}
	c.githubToken = os.Getenv("GITHUB_TOKEN")
	if os.Getenv("GITHUB_REPOS") != "" {
		c.githubRepos = strings.Split(os.Getenv("GITHUB_REPOS"), ",")
	}
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

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.githubToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	message := ""

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

		message += fmt.Sprintf("*%s*\n", repo)
		for _, pr := range pullRequests {
			timeAgo := humanize.Time(*pr.UpdatedAt)
			message += fmt.Sprintf(" • <%s|%s> - %s - updated %s\n", *pr.HTMLURL, *pr.Title, *pr.User.Login, timeAgo)
		}
		message += fmt.Sprintf("\n")
	}

	if message != "" {
		client := slack.New(conf.slackToken)

		opt := &slack.ChatPostMessageOpt{
			AsUser:    false,
			Username:  "purr",
			IconEmoji: ":purr:",
		}

		err := client.ChatPostMessage(conf.slackChannel, message, opt)
		if err != nil {
			fmt.Printf("while sending slack request: %s\n", err)
			os.Exit(1)
		}
	}
}
