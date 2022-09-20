package main

import (
	"fmt"
	"sync"

	"github.com/xanzy/go-gitlab"
)

func trawlGitLab(conf *Config, log Logger) <-chan *PullRequest {
	out := make(chan *PullRequest)

	// create a sync group that is used to close the out channel when all gitlab repos has been
	// trawled
	var wg sync.WaitGroup

	client, err := gitlab.NewClient(conf.GitLabToken, gitlab.WithBaseURL(conf.GitlabURL+"/api/v4"))
	if err != nil {
		usageAndExit(err.Error(), 1)
	}

	const status = "opened"

	// spin out each request to find PR on a repo into a separate goroutine
	for _, repo := range conf.GitLabRepos {

		// increment
		wg.Add(1)

		go func(repoName string) {
			defer wg.Done()
			log.Debugf("fetching GitLab PRs for %s\n", repoName)

			opts := &gitlab.ListProjectMergeRequestsOptions{
				State: gitlab.String(status),
			}
			pullRequests, _, err := client.MergeRequests.ListProjectMergeRequests(repoName, opts)
			if err != nil {
				log.Infof("Couldn't fetch PRs from GitLab (%s): %s\n", repoName, err)
				return
			}
			for _, pr := range pullRequests {
				out <- &PullRequest{
					ID:         pr.IID,
					Author:     pr.Author.Username,
					Assignee:   pr.Assignee.Username,
					Updated:    *pr.UpdatedAt,
					WebLink:    fmt.Sprintf("%s/%s/merge_requests/%d", conf.GitlabURL, repoName, pr.IID),
					Title:      pr.Title,
					Repository: repoName,
				}
			}
		}(repo)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
