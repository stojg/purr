# purr

[![CircleCI](https://circleci.com/gh/stojg/purr.svg?style=svg)](https://circleci.com/gh/stojg/purr)
[![Coverage Status](https://coveralls.io/repos/github/stojg/purr/badge.svg)](https://coveralls.io/github/stojg/purr)

Slack notifier for open pull requests.

## Screenshot

![example.png](./_docs/example.png)

## motivation

Pull requests are a good way of getting feedback from your team members. But pull requests can
easily get lost in the noise of daily communication and emails.

Slack can be quite easily be hooked up with GitHub/Gitlab to send notices when new pull requests are opened, but they
can be annoying. Purr can be configured to send a daily reminder, maybe in the morning with all the
pull requests that are waiting for review.

## features

- Get pull requests from Github
- Get merge requests from Gitlab
- Sends the summary to a slack channel
- Can be configured via a JSON file and environment variables
- Get all repositories for an Gitlab organisation
- Triggered via cron job or manually
- user configured filters

The slack message will be send by a user with the name `purr` and use the emoticon `:purr:` for as a slack icon.

## installation

If you are a gopher, you can install it via the usual:

`go get -u github.com/stojg/purr`

Compiled binaries are also available at [https://github.com/stojg/purr/releases](https://github.com/stojg/purr/releases)

## development

This project is using [glide](https://github.com/Masterminds/glide) for dependency management.

## configuration

purr can be configured with a JSON file and ENV variables. The ENV variables takes
precedence over the file configuration.

You will at least need a GitHub access token https://help.github.com/articles/creating-an-access-token-for-command-line-use/
with `repo` access and a Slack bot token https://api.slack.com/tokens.

Example JSON

```
{
  "github_token": "secret_token",
  "github_organisations": [
    "facebook"
  ],
  "github_users": [
    "stojg"
  ],
  "github_repos": [
    "user1/repo1",
    "user2/repo1"
  ],
  "gitlab_token": "secret_token",
  "gitlab_repos": [
    "project1/repo1",
    "project2/repo1"
  ],
  "gitlab_url": "https://www.example.com",
  "slack_token": "secret_token",
  "slack_channel": "myteamchat",
  "filters": {
    "wip": true,
    "users": [],
    "review": false
  }
}
```

Note that `github_organisations` will get all public and private repos and that `github_user` will only get the public
repos for a user due to how gitlab works.

The ENV variables are

```
export GITHUB_TOKEN="<super_secret_github token>"
export GITHUB_ORGANISATIONS - "facebook,twitter"
export GITHUB_USERS - "stojg,KentBeck"
export GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
export GITLAB_TOKEN="<super_secret_github token>"
export GITLAB_URL="http://example.com"
export GITLAB_REPOS="project1/repo1,project2/repo1"
export SLACK_TOKEN="<super_secret_slack_token>"
export SLACK_CHANNEL="my_slack_room"
export FILTER_USERS="user1,user2"
```

### filters

This is a description of the available filters and how to configure them:

###### wip bool, default: enabled

Will filter all requests which title begins with `WIP` or `[WIP]`, case-sensitive. If running on Github, will also filter out [Draft Pull Requests](https://github.blog/2019-02-14-introducing-draft-pull-requests/).

###### review bool, default: enabled

Will filter all Github request that has a "Changes requested" peer review

######  users list of strings, default: disabled

Will filter all pull requests where the author or assignee is not in the list of users

## run it

`purr --config my_team.json`

This is a one shot action, so you might want to put into a cron or a [systemd timer unit](https://wiki.archlinux.org/index.php/Systemd/Timers)

example `/etc/cron.d/purr` cron that runs purr 8am every day:

```
GITHUB_TOKEN="<super_secret_github token>"
SLACK_TOKEN="<super_secret_slack_token>"
GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
SLACK_CHANNEL="my_slack_room"
0 8 * * * username /usr/bin/purr --config /etc/purr/my_team.json
```
