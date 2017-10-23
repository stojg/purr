package main

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := newConfig("testdata/test_config.json")
	if err != nil {
		t.Error(err)
		return
	}

	validationErrors := config.validate()
	if len(validationErrors) != 0 {
		for _, err := range validationErrors {
			t.Errorf("Did not expect validation error: %+v", err)
		}
		return
	}

	if config.GitHubToken != "secret_github_token" {
		t.Errorf("expected GitHubToken to be 'secret_github_token', got '%s'", config.GitHubToken)
	}

	if len(config.GitHubOrganisations) != 1 {
		t.Errorf("Expected 1 GitHubOrganisations, got %d", len(config.GitHubRepos))
		return
	}
	if config.GitHubOrganisations[0] != "facebook" {
		t.Errorf("Expected GitHubOrganisationsto be '%s', got '%s'", "facebook", config.GitHubOrganisations[0])
		return
	}

	if len(config.GitHubUsers) != 1 {
		t.Errorf("Expected 1 GitHubUsers, got %d", len(config.GitHubUsers))
		return
	}
	if config.GitHubUsers[0] != "stojg" {
		t.Errorf("Expected first GitHubUsers to be '%s', got '%s'", "stojg", config.GitHubUsers[0])
		return
	}

	if len(config.GitHubRepos) != 2 {
		t.Errorf("Expected 2 GitHubRepos, got %d", len(config.GitHubRepos))
		return
	}
	if config.GitHubRepos[0] != "user1/repo1" {
		t.Errorf("Expected first GitHubRepos to be '%s', got '%s'", "user1/repo1", config.GitHubRepos[0])
		return
	}
	if config.GitHubRepos[1] != "user2/repo1" {
		t.Errorf("Expected second GitHubRepos to be '%s', got '%s'", "user2/repo1", config.GitHubRepos[1])
		return
	}

	if config.GitlabURL != "https://www.example.com" {
		t.Errorf("expected GitlabURL to be 'https://www.example.com', got '%s'", config.GitlabURL)
	}

	if config.GitLabToken != "secret_gitlab_token" {
		t.Errorf("expected GitLabToken to be 'secret_gitlab_token', got '%s'", config.GitLabToken)
	}

	if len(config.GitLabRepos) != 2 {
		t.Errorf("Expected 2 GitLabRepos, got %d", len(config.GitLabRepos))
		return
	}
	if config.GitLabRepos[0] != "project1/repo1" {
		t.Errorf("Expected first GitLabRepos to be '%s', got '%s'", "project1/repo1", config.GitLabRepos[0])
		return
	}
	if config.GitLabRepos[1] != "project2/repo1" {
		t.Errorf("Expected second GitLabRepos to be '%s', got '%s'", "project2/repo1", config.GitLabRepos[1])
		return
	}

	if config.SlackChannel != "myteamchat" {
		t.Errorf("Expected second SlackChannel to be 'myteamchat', got '%s'", config.SlackChannel)
		return
	}

	if config.SlackToken != "secret_slack_token" {
		t.Errorf("Expected second SlackToken to be 'secret_slack_token', got '%s'", config.SlackToken)
		return
	}

	if len(config.Filters.Users) != 2 {
		t.Errorf("Expected 2 filtered user, got got '%d'", len(config.Filters.Users))
		return
	}

	if !config.Filters.RequiresChanges {
		t.Errorf("Didnt expected filter 'review' to be disabled")
	}

	if !config.Filters.WorkInProgress {
		t.Errorf("Didnt expected filter 'wip' to be enabled")
	}

}

func TestNewConfig_NoFilters(t *testing.T) {
	config, err := newConfig("testdata/test_config_no_filters.json")
	if err != nil {
		t.Error(err)
		return
	}

	validationErrors := config.validate()
	if len(validationErrors) != 0 {
		for _, err := range validationErrors {
			t.Errorf("Did not expect validation error: %+v", err)
		}
		return
	}

	if len(config.Filters.Users) != 0 {
		t.Errorf("Expected 0 filtered users, got got '%d'", len(config.Filters.Users))
		return
	}

	if config.Filters.RequiresChanges {
		t.Errorf("Didnt expected filter 'review' to be enabled")
	}

	if config.Filters.WorkInProgress {
		t.Errorf("Didnt expected filter 'wip' to be enabled")
	}

}
