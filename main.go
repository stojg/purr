package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"sync"
	"time"
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
	cliOutput  bool
)

func init() {
	flag.StringVar(&configFile, "config", "", "Read config from FILE")
	flag.BoolVar(&debug, "d", false, "run in debug mode")
	flag.BoolVar(&cliOutput, "o", false, "output to CLI rather than slack")

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

	// these function will return channels that will emit a list of pull requests
	// on channels and close the channel when they are done
	gitHubPRs := trawlGitHub(conf)
	gitLabPRs := trawlGitLab(conf)

	// Merge the in channels into of channel and close it when the inputs are done
	prs := merge(gitHubPRs, gitLabPRs)

	// filter out pull requests that we don't want to send
	filteredPRs := filter(conf, prs)

	// format takes a channel of pull requests and returns a message that groups
	// pull request into repos and formats them into a slack friendly format
	message := format(filteredPRs)

	// Output what slack will send if we are in debug mode
	if debug {
		logrus.Debugf("Final message:\n%s", message)
	}

	// If -o is set, just output
	if cliOutput {
		fmt.Print(message)

		// Otherwise send to slack
	} else {
		postToSlack(conf, message)
	}
}

func trawlGitHub(conf *Config) <-chan *PullRequest {

	out := make(chan *PullRequest)

	// create a sync group that is used to close the out channel when all github repos has been
	// trawled
	var wg sync.WaitGroup

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.GitHubToken})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// check for wildcards in the repo name and expand them into individual repos
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
		allRepos, _, err := client.Repositories.List(context.Background(), repoParts[0], nil)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for i := range allRepos {
			repos = append(repos, fmt.Sprintf("%s/%s", repoParts[0], *allRepos[i].Name))
		}
	}

	// spin out each request to find PRs on a repo into a separate goroutine so we fetch them
	// asynchronous
	for _, repo := range repos {

		// increment the wait group
		wg.Add(1)

		go func(repoName string) {
			// when finished, decrement the wait group
			defer wg.Done()
			logrus.Debugf("Starting fetch from %s", repoName)

			parts := strings.Split(repoName, "/")

			// nextPage keeps track of of the current page of the paginataed response from the
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
				pullRequests, resp, err := client.PullRequests.List(context.Background(), parts[0], parts[1], options)
				if err != nil {
					logrus.Errorf("While fetching PRs from GitHub (%s/%s): %s", parts[0], parts[1], err)
					return
				}

				// transform the GitHub pull request struct into a provider agnostic struct
				for _, pr := range pullRequests {
					pullRequest := &PullRequest{
						ID:         *pr.Number,
						Author:     *pr.User.Login,
						Updated:    *pr.UpdatedAt,
						WebLink:    *pr.HTMLURL,
						Title:      *pr.Title,
						Repository: fmt.Sprintf("%s/%s", parts[0], parts[1]),
					}
					if pr.Assignee != nil {
						pullRequest.Assignee = *pr.Assignee.Login
					}

					// push to the outchannel
					out <- pullRequest
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
		logrus.Debugf("Done with github")
		close(out)
	}()

	return out
}

func trawlGitLab(conf *Config) <-chan *PullRequest {
	out := make(chan *PullRequest)

	// create a sync group that is used to close the out channel when all gitlab repos has been
	// trawled
	var wg sync.WaitGroup

	client := gitlab.NewClient(nil, conf.GitLabToken)
	if err := client.SetBaseURL(conf.GitlabURL + "/api/v3"); err != nil {
		usageAndExit(err.Error(), 1)
	}

	status := "opened"
	options := &gitlab.ListMergeRequestsOptions{State: &status}

	// spin out each request to find PR on a repo into a separate goroutine
	for _, repo := range conf.GitLabRepos {

		// increment
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
		logrus.Debugf("Done with gitlab")
		close(out)
	}()

	return out
}

// merge merges several channels into one output channel (fan-in)
func merge(channels ...<-chan *PullRequest) <-chan *PullRequest {
	var wg sync.WaitGroup
	out := make(chan *PullRequest)

	// Start an output goroutine for each input channel in channels. output copies values from prs
	// to out until prs is closed, then calls wg.Done
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

	// Start a goroutine to close out once all the output goroutines are done. This must start after
	// the wg.Add call
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// filter removes pull requests that should not show up in the final message, this could
// include PRs marked as Work in Progress or where users are not in the whitelist
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
	numPRs := 0
	var oldest *PullRequest
	lastUpdated := time.Now()
	for pr := range prs {
		if pr.Updated.Before(lastUpdated) {
			oldest = pr
			lastUpdated = pr.Updated
		}
		numPRs++
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

	if numPRs > 0 {
		fmt.Fprintf(buf, "\nThere are currently %d open pull requests", numPRs)
		fmt.Fprintf(buf, " and the oldest (<%s|PR #%d>) was updated %s\n", oldest.WebLink, oldest.ID, humanize.Time(oldest.Updated))
	}
	return buf
}

// postToSlack will post the message to Slack. It will divide the message into smaller message if
// it's more than 30 lines long due to a max message size limitation enforced by the Slack API
func postToSlack(conf *Config, message fmt.Stringer) {

	const maxLines = 30

	if message.String() == "" {
		return
	}

	client := slack.New(conf.SlackToken)
	opt := &slack.ChatPostMessageOpt{
		AsUser:    false,
		Username:  "purr",
		IconEmoji: ":purr:",
	}

	// Don't send to large messages, send a new message per 40 new lines
	lines := strings.Split(message.String(), "\n")
	lineBuffer := make([]string, maxLines)
	for i := range lines {
		lineBuffer = append(lineBuffer, lines[i])
		if len(lineBuffer) == cap(lineBuffer) || i+1 == len(lines) {
			msg := strings.Join(lineBuffer, "\n")
			if msg != "" {
				if err := client.ChatPostMessage(conf.SlackChannel, msg, opt); err != nil {
					logrus.Errorf("Slack: %s", err)
					os.Exit(1)
				}
			}
			lineBuffer = make([]string, cap(lineBuffer))
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
