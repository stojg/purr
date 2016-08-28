package main

import (
	"fmt"
	"github.com/bluele/slack"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

type config struct {
	repos []string
}

func main() {

	envRepos := os.Getenv("GITHUB_REPOS")
	if envRepos == "" {
		panic("No repos defined in ENV 'GITHUB_REPOS'")
	}

	repos := strings.Split(envRepos, ",")

	if len(repos) == 0 {
		panic("No repos defined in ENV 'GITHUB_REPOS'")
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		panic("No token defined in ENV 'GITHUB_TOKEN'")
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		panic("No token found in ENV 'SLACK_TOKEN'")
	}

	slackChannel := os.Getenv("SLACK_CHANNEL")
	if slackChannel == "" {
		panic("No slack channel found in ENV 'SLACK_CHANNEL'")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	message := ""

	for _, gitRepo := range repos {

		parts := strings.Split(gitRepo, "/")

		prs, _, err := client.PullRequests.List(parts[0], parts[1], nil)
		if err != nil {
			panic(err)
		}
		if len(prs) == 0 {
			continue
		}
		message += fmt.Sprintf("*%s*\n", gitRepo)

		for _, pr := range prs {
			message += fmt.Sprintf(" • <%s|%s> - %s\n", *pr.HTMLURL, *pr.Title, *pr.User.Login)
		}
		message += fmt.Sprintf("\n")
	}

	if message != "" {
		slack := slack.New(slackToken)
		err := slack.ChatPostMessage(slackChannel, message, nil)
		if err != nil {
			panic(err)
		}
	}
}
