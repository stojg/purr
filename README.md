# purr

Send a list of open github pull requests to a Slack channel.

## Motivation

Peer reviews are a really good way to work and communicate as a team. Unfortunately, it's easy for open pull requests to be to get lost in the daily noise of emails, meetings other direct messaging.

Two of the corner stones of agile development is to shorten the feedback loop and reduce waste (waste as in work done, but not deployed). Open pull requests are a good example of waiting for feedback and waste. They should either be worked on until they can be merged and deployed or closed.

This little program is meant to be used as a periodic reminder of Github and Gitlab pull requests that are open.

It will ignore Pull Requests that contains `[WIP]` or `WIP:` in the title.

## installation

`go get -u github.com/stojg/purr`

... or download one of the binaries for your platform from https://github.com/stojg/purr/releases

If you want to cross-compile for other platform, see http://golangcookbook.com/chapters/running/cross-compiling/

## configuration

```
export GITHUB_TOKEN="<super_secret_github token>"
export GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
export SLACK_TOKEN="<super_secret_slack_token>"
export SLACK_CHANNEL="my_slack_room"
```

Optional configuration for Gitlab is

```
export GITLAB_TOKEN="<super_secret_gitlab token>"
export GITLAB_URL="https://gitlab.example.com"
export GITLAB_REPOS="namespace/project1,namespace/project2"
```

It will send the message as a bot user with the name `purr` and use the emoticon
`:purr:` as its avatar in slack.

## run it

`purr`

This is a one shot action, so you might want to put into a cron or a [systemd timer unit](https://wiki.archlinux.org/index.php/Systemd/Timers)

example `/etc/cron.d/purr` cron that runs purr 8am every day:

```
GITHUB_TOKEN="<super_secret_github token>"
SLACK_TOKEN="<super_secret_slack_token>"
GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
SLACK_CHANNEL="my_slack_room"
0 8 * * * username /usr/bin/purr
```

## End result:

![example.png](./_docs/example.png)
