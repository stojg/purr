package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Config contains the settings from the user
type Config struct {
	GitHubToken   string   `json:"github_token"`
	GitHubRepos   []string `json:"github_repos"`
	GitLabToken   string   `json:"gitlab_token"`
	GitLabRepos   []string `json:"gitlab_repos"`
	GitlabURL     string   `json:"gitlab_url"`
	SlackToken    string   `json:"slack_token"`
	SlackChannel  string   `json:"slack_channel"`
	UserWhiteList []string `json:"user_whitelist"`
}

func newConfig() *Config {

	c := &Config{}

	if configFile != "" {
		file, e := ioutil.ReadFile(configFile)
		if e != nil {
			usageAndExit(e.Error(), 1)
		}

		if err := json.Unmarshal(file, &c); err != nil {
			usageAndExit(err.Error(), 1)
		}
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		c.GitHubToken = os.Getenv("GITHUB_TOKEN")
	}
	if os.Getenv("GITHUB_REPOS") != "" {
		c.GitHubRepos = strings.Split(os.Getenv("GITHUB_REPOS"), ",")
	}
	if os.Getenv("GITLAB_TOKEN") != "" {
		c.GitLabToken = os.Getenv("GITLAB_TOKEN")
	}
	if os.Getenv("GITLAB_REPOS") != "" {
		c.GitLabRepos = strings.Split(os.Getenv("GITLAB_REPOS"), ",")
	}
	if os.Getenv("GITLAB_URL") != "" {
		c.GitlabURL = os.Getenv("GITLAB_URL")
	}
	if os.Getenv("SLACK_TOKEN") != "" {
		c.SlackToken = os.Getenv("SLACK_TOKEN")
	}
	if os.Getenv("SLACK_CHANNEL") != "" {
		c.SlackChannel = os.Getenv("SLACK_CHANNEL")
	}
	if os.Getenv("USER_WHITELIST") != "" {
		c.UserWhiteList = strings.Split(os.Getenv("USER_WHITELIST"), ",")
	}
	return c
}

func (c *Config) validate() []error {
	var errors []error
	if c.GitHubToken == "" {
		errors = append(errors, fmt.Errorf("GitHub token cannot be empty"))
	}
	if len(c.GitHubRepos) == 0 {
		errors = append(errors, fmt.Errorf("GitHub repos cannot be empty"))
	}
	if c.SlackToken == "" {
		errors = append(errors, fmt.Errorf("Slack token cannot be empty"))
	}
	if c.SlackChannel == "" {
		errors = append(errors, fmt.Errorf("Slack channel cannot be empty"))
	}

	return errors
}

func configHelp() {
	fmt.Fprintln(os.Stderr, "\npurr requrires configuration to be either in a config file or set the ENV")

	fmt.Fprintln(os.Stderr, "\nThe configuration file (--config) looks like this:")

	exampleConfig := &Config{
		GitHubToken:   "secret_token",
		GitHubRepos:   []string{"user1/repo1", "user2/repo1"},
		GitLabToken:   "secret_token",
		GitLabRepos:   []string{"project1/repo1", "project2/repo1"},
		GitlabURL:     "https://www.example.com",
		SlackToken:    "secret_token",
		SlackChannel:  "myteamchat",
		UserWhiteList: []string{"user1", "user2"},
	}
	b, err := json.MarshalIndent(exampleConfig, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	fmt.Fprintf(os.Stderr, "\n%s\n\n", b)

	fmt.Fprint(os.Stderr, "The above configuration can be overridden with ENV variables:\n\n")
	fmt.Fprintln(os.Stderr, " * GITHUB_TOKEN")
	fmt.Fprintln(os.Stderr, " * GITHUB_REPOS - comma separated list")
	fmt.Fprintln(os.Stderr, " * GITLAB_TOKEN")
	fmt.Fprintln(os.Stderr, " * GITLAB_URL")
	fmt.Fprintln(os.Stderr, " * GITLAB_REPOS - comma separated list")
	fmt.Fprintln(os.Stderr, " * SLACK_TOKEN")
	fmt.Fprintln(os.Stderr, " * SLACK_CHANNEL")
	fmt.Fprintln(os.Stderr, " * USER_WHITELIST - comma separated list")
}
