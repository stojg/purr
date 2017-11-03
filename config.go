package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//"github.com/mitchellh/mapstructure"
)

// Config contains the settings from the user
type Config struct {
	GitHubToken         string   `json:"github_token"`
	GitHubOrganisations []string `json:"github_organisations"`
	GitHubUsers         []string `json:"github_users"`
	GitHubRepos         []string `json:"github_repos"`
	GitLabToken         string   `json:"gitlab_token"`
	GitLabRepos         []string `json:"gitlab_repos"`
	GitlabURL           string   `json:"gitlab_url"`
	SlackToken          string   `json:"slack_token"`
	SlackChannel        string   `json:"slack_channel"`
	Filters             *Filters `json:"filters"`
}

func newConfig(filePath string) (*Config, error) {

	config := &Config{
		Filters: &Filters{},
	}

	// the config Filters is an slice of interfaces, so we need to manually set defaults and add them to the Config
	filters := struct {
		Users  UserFilter
		WIP    WIPFilter    `json:"wip"`
		Review ReviewFilter `json:"review"`
	}{
		WIP: true,
		Review: true,
	}

	if filePath != "" {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			return config, fmt.Errorf("Error during config read: %s", err)
		}

		// populate as much as possible into the config
		if err := json.Unmarshal(file, &config); err != nil {
			return config, fmt.Errorf("Error during config read: %s", err)
		}

		if err := json.Unmarshal(file, &filters); err != nil {
			return config, fmt.Errorf("Error during config read: %s", err)
		}
	}

	if os.Getenv("GITHUB_TOKEN") != "" {
		config.GitHubToken = os.Getenv("GITHUB_TOKEN")
	}
	if os.Getenv("GITHUB_ORGANISATIONS") != "" {
		config.GitHubOrganisations = strings.Split(os.Getenv("GITHUB_ORGANISATIONS"), ",")
	}
	if os.Getenv("GITHUB_USERS") != "" {
		config.GitHubUsers = strings.Split(os.Getenv("GITHUB_USERS"), ",")
	}
	if os.Getenv("GITHUB_REPOS") != "" {
		config.GitHubRepos = strings.Split(os.Getenv("GITHUB_REPOS"), ",")
	}
	if os.Getenv("GITLAB_TOKEN") != "" {
		config.GitLabToken = os.Getenv("GITLAB_TOKEN")
	}
	if os.Getenv("GITLAB_REPOS") != "" {
		config.GitLabRepos = strings.Split(os.Getenv("GITLAB_REPOS"), ",")
	}
	if os.Getenv("GITLAB_URL") != "" {
		config.GitlabURL = os.Getenv("GITLAB_URL")
	}
	if os.Getenv("SLACK_TOKEN") != "" {
		config.SlackToken = os.Getenv("SLACK_TOKEN")
	}
	if os.Getenv("SLACK_CHANNEL") != "" {
		config.SlackChannel = os.Getenv("SLACK_CHANNEL")
	}
	if os.Getenv("FILTER_USERS") != "" {
		filters.Users = strings.Split(os.Getenv("FILTER_USERS"), ",")
	}
	if os.Getenv("FILTER_WIP") != "" {
		filters.WIP = os.Getenv("FILTER_WIP") == "true"
	}
	if os.Getenv("FILTER_REVIEW") != "" {
		filters.Review = os.Getenv("FILTER_REVIEW") == "true"
	}

	config.Filters.Add(filters.Users)
	config.Filters.Add(filters.Review)
	config.Filters.Add(filters.WIP)

	return config, nil
}

func (c *Config) validate() []error {
	var errors []error
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
		GitHubToken:         "secret_token",
		GitHubOrganisations: []string{"facebook"},
		GitHubUsers:         []string{"stojg"},
		GitHubRepos:         []string{"user1/repo1", "user2/repo1"},
		GitLabToken:         "secret_token",
		GitLabRepos:         []string{"project1/repo1", "project2/repo1"},
		GitlabURL:           "https://www.example.com",
		SlackToken:          "secret_token",
		SlackChannel:        "myteamchat",
		Filters:             &Filters{},
	}

	b, err := json.MarshalIndent(exampleConfig, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	fmt.Fprintf(os.Stderr, "\n%s\n\n", b)

	fmt.Fprint(os.Stderr, "The above configuration can be overridden with ENV variables:\n\n")
	fmt.Fprintln(os.Stderr, " * GITHUB_TOKEN")
	fmt.Fprintln(os.Stderr, " * GITHUB_ORGANISATIONS - comma separated list")
	fmt.Fprintln(os.Stderr, " * GITHUB_USERS - comma separated list")
	fmt.Fprintln(os.Stderr, " * GITHUB_REPOS - comma separated list")
	fmt.Fprintln(os.Stderr, " * GITLAB_TOKEN")
	fmt.Fprintln(os.Stderr, " * GITLAB_URL")
	fmt.Fprintln(os.Stderr, " * GITLAB_REPOS - comma separated list")
	fmt.Fprintln(os.Stderr, " * SLACK_TOKEN")
	fmt.Fprintln(os.Stderr, " * SLACK_CHANNEL")
	fmt.Fprintln(os.Stderr, " * FILTER_USERS - comma separated list")
	fmt.Fprintln(os.Stderr, " * FILTER_WIP - 'true' or 'false'")
	fmt.Fprintln(os.Stderr, " * FILTER_REVIEW - 'true' or 'false'")
}
