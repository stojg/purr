package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func trawlGitHub(conf *Config, log Logger) <-chan *PullRequest {

	out := make(chan *PullRequest)

	// create a sync group that is used to close the out channel when all github repos has been
	// trawled
	var wg sync.WaitGroup

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GitHubToken})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	var repos []string

	// check for a organisation and all it's repositories
	for _, organisationName := range conf.GitHubOrganisations {
		// first try listing by organisation
		allRepos, _, err := client.Repositories.ListByOrg(context.Background(), organisationName, nil)
		if err != nil {
			log.Infof("Failed getting repositories for GitHub organisation %s: %v\n", organisationName, err)
			continue
		}
		for i := range allRepos {
			repos = append(repos, *allRepos[i].FullName)
		}
	}

	for _, user := range conf.GitHubUsers {
		// first try listing by organisation
		allRepos, _, err := client.Repositories.List(context.Background(), user, nil)
		if err != nil {
			log.Infof("Failed getting repositories for GitHub user %s: %v\n", user, err)
			continue
		}
		for i := range allRepos {
			repos = append(repos, *allRepos[i].FullName)
		}
	}

	for _, repoName := range conf.GitHubRepos {
		repoParts := strings.Split(repoName, "/")
		if len(repoParts) != 2 {
			log.Infof("%s is not a valid GitHub repository\n", repoName)
			continue
		}
		repos = append(repos, repoName)
	}

	// spin out each request to find PRs on a repo into a separate goroutine so we fetch them
	// asynchronous
	for _, repo := range repos {

		// increment the wait group
		wg.Add(1)

		go func(repoName string) {
			// when finished, decrement the wait group
			defer wg.Done()

			parts := strings.Split(repoName, "/")

			// nextPage keeps track of of the current page of the paginated response from the
			// GitHub API
			nextPage := 1
			for {
				// options for the request for PRs
				options := &github.PullRequestListOptions{
					State:     "open",
					Sort:      "updated",
					Direction: "desc",
					ListOptions: github.ListOptions{
						Page: nextPage,
					},
				}

				// get the pull requests
				log.Debugf("fetching all PRs for GitHub repo %s\n", repoName)
				pullRequests, resp, err := client.PullRequests.List(context.Background(), parts[0], parts[1], options)
				if err != nil {
					log.Infof("couldn't fetch PRs from GitHub (%s): %s\n", repoName, err)
					return
				}

				// transform the GitHub pull request struct into a provider agnostic struct
				for _, pr := range pullRequests {
					wg.Add(1)

					// Get the github reviews and push result onto out when done
					go func(pr *github.PullRequest) {
						defer wg.Done()

						requiresChanges, approved := trawlGitHubReviews(client, parts[0], parts[1], *pr.Number, log)

						pullRequest := &PullRequest{
							ID:              *pr.Number,
							Author:          *pr.User.Login,
							Updated:         *pr.UpdatedAt,
							WebLink:         *pr.HTMLURL,
							Title:           *pr.Title,
							RequiresChanges: requiresChanges,
							Approved:        approved,
							Repository:      fmt.Sprintf("%s/%s", parts[0], parts[1]),
						}
						if pr.Assignee != nil {
							pullRequest.Assignee = *pr.Assignee.Login
						}
						out <- pullRequest
					}(pr)
				}

				// the GitHub API returns 0 as the LastPage if there are no more pages of result
				if resp.LastPage == 0 {
					break
				}
				nextPage++

			}
		}(repo)
	}

	// Spin off a go routine that will close the channel when all repos have finished
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// trawlGitHubReviews goes through the reviews of a single PR and returns a few flags: requiresChanges, approved
func trawlGitHubReviews(client *github.Client, owner string, repo string, number int, log Logger) (bool, bool) {
	requiresChanges := false
	approved := false

	nextPage := 1
	for {
		options := &github.ListOptions{
			Page: nextPage,
		}

		// get the reviews for the PR
		pullRequestReviews, resp, err := client.PullRequests.ListReviews(context.Background(), owner, repo, number, options)
		if err != nil {
			log.Infof("Couldn't fetch PR reviews from GitHub (%s/%s#%d): %s\n", owner, repo, number, err)
			return false, false
		}

		// the list of reviews is in chronological order, which means that if a review requires changes
		// after it's been approved, the PRs approval state is false
		for _, review := range pullRequestReviews {
			if *review.State == "CHANGES_REQUESTED" {
				requiresChanges = true
				approved = false
			}
			if *review.State == "APPROVED" {
				approved = true
				requiresChanges = false
			}
		}

		// the GitHub API returns 0 as the LastPage if there are no more pages of result
		if resp.LastPage == 0 {
			break
		}
		nextPage++
	}

	return requiresChanges, approved
}
