package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bluele/slack"
	"github.com/dustin/go-humanize"
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
}

func main() {

	logger := NewStdOutLogger(debug)

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
	gitHubPRs := trawlGitHub(conf, logger)
	gitLabPRs := trawlGitLab(conf, logger)

	// Merge the in channels into of channel and close it when the inputs are done
	prs := merge(gitHubPRs, gitLabPRs)

	// filter out pull requests that we don't want to send
	filteredPRs := filter(conf.Filters, prs, logger)

	// format takes a channel of pull requests and returns a message that groups
	// pull request into repos and formats them into a slack friendly format
	message := format(conf.Filters, filteredPRs)

	if message.String() == "" {
		logger.Debugf("No PRs found\n")
	} else if cliOutput {
		fmt.Print(message)
	} else {
		err := postToSlack(conf, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not send to slack: %v\n", err)
			os.Exit(1)
		}
	}
}

// merge merges several channels into one output channel (fan-in)
func merge(channels ...<-chan *PullRequest) <-chan *PullRequest {
	out := make(chan *PullRequest)

	var wg sync.WaitGroup
	wg.Add(len(channels))

	// merges all in channels into an out channel
	for _, c := range channels {
		go func(prs <-chan *PullRequest) {
			for pr := range prs {
				out <- pr
			}
			wg.Done()
		}(c)
	}

	// we have to wait until all input channels are closed before closing the out channel
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// filter removes pull requests that should not show up in the final message, this could
// include PRs marked as Work in Progress or where users are not in the whitelist
func filter(filters *Filters, in <-chan *PullRequest, log Logger) chan *PullRequest {
	out := make(chan *PullRequest)

	go func() {
		for pr := range in {
			if filters.Filter(pr) {
				out <- pr
			} else {
				log.Debugf("filtered PR '%s' (%s) \n", pr.Title, pr.WebLink)
			}
		}
		close(out)
	}()
	return out
}

// format converts all pull requests into a message that is grouped by repo formatted for slack
func format(filters *Filters, prs <-chan *PullRequest) fmt.Stringer {
	var numPRs int
	var oldest *PullRequest
	repositories := make(map[string][]*PullRequest)
	lastUpdated := time.Now()

	// loop through all PRs, will stop when the channel is closed
	for pr := range prs {
		// update the oldest pull request
		if pr.Updated.Before(lastUpdated) {
			oldest = pr
			lastUpdated = pr.Updated
		}
		// update the total count of PRs
		numPRs++
		// group PRs with their repository
		if _, ok := repositories[pr.Repository]; !ok {
			repositories[pr.Repository] = make([]*PullRequest, 0)
		}
		repositories[pr.Repository] = append(repositories[pr.Repository], pr)
	}

	// create the slack message in "slack" format
	buf := &bytes.Buffer{}
	for repo, prs := range repositories {
		fmt.Fprintf(buf, "*%s*\n", repo)
		for i := range prs {
			fmt.Fprintf(buf, "%s\n", prs[i])
		}
		fmt.Fprint(buf, "\n")
	}

	// summary
	if numPRs > 0 {
		fmt.Fprintf(buf, "\nThere are currently %d open pull requests", numPRs)
		fmt.Fprintf(buf, " and the oldest (<%s|PR #%d>) was updated %s\n", oldest.WebLink, oldest.ID, humanize.Time(oldest.Updated))
	}
	fmt.Fprintf(buf, "%d pull requests was filtered from this result\n", filters.NumFiltered())
	return buf
}

// postToSlack will post the message to Slack. It will divide the message into smaller message if
// it's more than 30 lines long due to a max message size limitation enforced by the Slack API
func postToSlack(conf *Config, message fmt.Stringer) error {

	const maxLines = 30

	client := slack.New(conf.SlackToken)
	opt := &slack.ChatPostMessageOpt{
		AsUser:    false,
		Username:  "purr",
		IconEmoji: ":purr:",
	}

	// Don't send too large messages, send a new message per maxLines new lines
	lines := strings.Split(message.String(), "\n")
	lineBuffer := make([]string, maxLines)
	for i := range lines {
		lineBuffer = append(lineBuffer, lines[i])
		if len(lineBuffer) == cap(lineBuffer) || i+1 == len(lines) {
			msg := strings.Join(lineBuffer, "\n")
			if msg != "" {
				if err := client.ChatPostMessage(conf.SlackChannel, msg, opt); err != nil {
					return err
				}
			}
			lineBuffer = make([]string, cap(lineBuffer))
		}
	}
	return nil
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
