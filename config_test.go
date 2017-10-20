package main

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := newConfig("testdata/test_config.json")
	if err != nil {
		t.Error(err)
		return
	}

	if config.GitHubToken != "super_secret" {
		t.Errorf("expected GitHubToken to be 'super_secret', got '%s'", config.GitHubToken)
	}

	if len(config.GitHubOrganisations) != 1 {
		t.Errorf("Expected 1 GitHubOrganisations, got %d", len(config.GitHubRepos))
		return
	}

	if config.GitHubOrganisations[0] != "twitter" {
		t.Errorf("Expected GitHubOrganisationsto be '%s', got '%s'", "twitter", config.GitHubOrganisations[0])
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

	if config.GitHubRepos[0] != "facebook/emitter" {
		t.Errorf("Expected first GitHubRepos to be '%s', got '%s'", "facebook/emitter", config.GitHubRepos[0])
		return
	}
	if config.GitHubRepos[1] != "facebook/buck" {
		t.Errorf("Expected second GitHubRepos to be '%s', got '%s'", "facebook/buck", config.GitHubRepos[1])
		return
	}
}
