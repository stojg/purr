package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	//"github.com/xanzy/go-gitlab"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"sync"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "purr - %s\n"
	// VERSION is the binary version.
	VERSION = "v0.4.0"
)

var (
	configFile string
	debug      bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Read config from FILE")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
		configHelp()
	}
	flag.Parse()

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {

	conf, err := newConfig(configFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if errors := conf.validate(); len(errors) > 0 {
		buf := &bytes.Buffer{}
		for i := range errors {
			fmt.Fprintln(buf, errors[i].Error())
		}
		usageAndExit(buf.String(), 1)
	}

	gitHubPRs := trawlGitHub(conf)
	gitLabPRs := trawlGitLab(conf)

	prs := merge(gitHubPRs, gitLabPRs)

	filteredPRs := filter(conf, prs)

	message := format(filteredPRs)

	if debug {
		logrus.Debugf("Final message:\n%s", message)
	}

	postToSlack(conf, message)
}

func trawlGitHub(conf *Config) <-chan *PullRequest {

	out := make(chan *PullRequest)

	var wg sync.WaitGroup

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GitHubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// check for wildcards and expand them
	var repos []string
	for _, repoName := range conf.GitHubRepos {
		repoParts := strings.Split(repoName, "/")
		if len(repoParts) != 2 {
			logrus.Errorf("%s is not a valid GitHub repository\n", repoName)
			continue
		}
		if repoParts[1] != "*" {
			repos = append(repos, repoName)
			continue
		}
		logrus.Debugf("expanding wildcard on %s", repoName)
		allRepos, _, err := client.Repositories.List(repoParts[0], nil)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for i := range allRepos {
			repos = append(repos, fmt.Sprintf("%s/%s", repoParts[0], *allRepos[i].Name))
		}
	}

	// spin out each request to find PR on a repo into a separate goroutine
	for _, repo := range repos {

		wg.Add(1)

		go func(repoName string) {
			defer wg.Done()
			logrus.Debugf("Starting fetch from %s", repoName)

			parts := strings.Split(repoName, "/")

			pullRequests, _, err := client.PullRequests.List(parts[0], parts[1], nil)
			if err != nil {
				logrus.Errorf("While fetching PRs from GitHub (%s/%s): %s", parts[0], parts[1], err)
				return
			}

			for _, pr := range pullRequests {
				pullRequest := &PullRequest{
					Author:     *pr.User.Login,
					TimeAgo:    humanize.Time(*pr.UpdatedAt),
					WebLink:    *pr.HTMLURL,
					Title:      *pr.Title,
					Repository: fmt.Sprintf("%s/%s", parts[0], parts[1]),
				}
				if pr.Assignee != nil {
					pullRequest.Assignee = *pr.Assignee.Login
				}
				out <- pullRequest
			}
		}(repo)
	}

	go func() {
		wg.Wait()
		logrus.Debugf("Done with github")
		close(out)
	}()

	return out
}

func trawlGitLab(conf *Config) <-chan *PullRequest {
	out := make(chan *PullRequest)

	var wg sync.WaitGroup

	client := gitlab.NewClient(nil, conf.GitLabToken)
	if err := client.SetBaseURL(conf.GitlabURL + "/api/v3"); err != nil {
		usageAndExit(err.Error(), 1)
	}

	status := "opened"
	options := &gitlab.ListMergeRequestsOptions{State: &status}

	// spin out each request to find PR on a repo into a separate goroutine
	for _, repo := range conf.GitLabRepos {

		wg.Add(1)

		go func(repoName string) {
			defer wg.Done()

			pullRequests, _, err := client.MergeRequests.ListMergeRequests(repoName, options)
			if err != nil {
				logrus.Errorf("While fetching PRs from GitLab (%s): %s", repoName, err)
				return
			}
			for _, pr := range pullRequests {
				out <- &PullRequest{
					Author:     pr.Author.Username,
					Assignee:   pr.Assignee.Username,
					TimeAgo:    humanize.Time(*pr.UpdatedAt),
					WebLink:    fmt.Sprintf("%s/%s/merge_requests/%d", conf.GitlabURL, repoName, pr.IID),
					Title:      pr.Title,
					Repository: repoName,
				}
			}
		}(repo)
	}

	go func() {
		wg.Wait()
		logrus.Debugf("Done with gitlab")
		close(out)
	}()

	return out
}

func filter(conf *Config, in <-chan *PullRequest) chan *PullRequest {
	out := make(chan *PullRequest)

	go func() {
		for list := range in {
			if !list.isWIP() && list.isWhiteListed(conf) {
				out <- list
			} else {
				logrus.Debugf("filtered pr '%s'", list.Title)
			}
		}
		close(out)
	}()

	return out
}

// format converts all pull requests into a message that is grouped by repo formatted for slack
func format(prs <-chan *PullRequest) fmt.Stringer {
	grouped := make(map[string][]*PullRequest)
	for pr := range prs {
		if _, ok := grouped[pr.Repository]; !ok {
			grouped[pr.Repository] = make([]*PullRequest, 0)
		}
		grouped[pr.Repository] = append(grouped[pr.Repository], pr)
	}

	buf := &bytes.Buffer{}
	for repo, prs := range grouped {
		fmt.Fprintf(buf, "*%s*\n", repo)
		for i := range prs {
			fmt.Fprintf(buf, "%s\n", prs[i])
		}
		fmt.Fprint(buf, "\n")
	}
	return buf
}

// merge merges several channels into one output channel (fan-in)
func merge(channels ...<-chan *PullRequest) <-chan *PullRequest {
	var wg sync.WaitGroup
	out := make(chan *PullRequest)

	// Start an output goroutine for each input channel in cs. output copies values from c to out
	// until c is closed, then calls wg.Done.
	output := func(prs <-chan *PullRequest) {
		for pr := range prs {
			out <- pr
		}
		wg.Done()
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.  This must start
	// after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func postToSlack(conf *Config, message fmt.Stringer) {
	if message.String() != "" {
		client := slack.New(conf.SlackToken)
		opt := &slack.ChatPostMessageOpt{
			AsUser:    false,
			Username:  "purr",
			IconEmoji: ":purr:",
		}
		err := client.ChatPostMessage(conf.SlackChannel, message.String(), opt)
		if err != nil {
			logrus.Errorf("Sending slack request: %s\n", err)
			os.Exit(1)
		}
	}
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
