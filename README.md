# purr

Send a list of open github pull requests to a Slack channel.


## Screenshot

![example.png](./_docs/example.png)

## motivation

Peer reviews are a really good way to work and communicate as a team. Unfortunately it's easy for
pull requests to be to get lost in the daily noise of emails, meetings and slack.

Two of the corner stones of agile development is to shorten the feedback loop and reduce waste
(waste as in work done, but not deployed). Open pull requests are a good example of waiting for
feedback and waste. They should either be worked on until they can be merged and deployed or closed.

## features

- Sends the summary to a slack channel
- List open pull requests from GitHub
- Can be configured via a JSON file and environment variables
- Use `user/*` in the repo configuration to get PRs for all repositories for that user
- Can ignore PRs if the author or assignee is not in a whitelist
- Ignores pull requests that contains `[WIP]` or `WIP:` in the title
- Triggered via cron job or manually


The slack message will be send by a user with the name `purr` and use the slack icon is the emoticon
`:purr:`.


Experimental features

- GitLab integration

## installation

If you are a gopher, you can install it via the usual:

`go get -u github.com/stojg/purr`

Compiled binaries are also available at https://github.com/stojg/purr/releases

## configuration

purr can be configured with a JSON file and ENV variables. The ENV variables takes
precedence over the file configuration.

You will at least need a GitHub access token https://help.github.com/articles/creating-an-access-token-for-command-line-use/
with `repo` access and a Slack bot token https://api.slack.com/tokens.

Example JSON

```
{
  "github_token": "secret_token",
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
  "user_whitelist": [
    "user1",
    "user2"
  ]
}
```

To get all repos for a github user / organisation, replace the repository name with a star:

```
"github_repos": [ "user1/*" ]

```


The ENV variables are

```
export GITHUB_TOKEN="<super_secret_github token>"
export GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
export GITLAB_TOKEN="<super_secret_github token>"
export GITLAB_URL="http://example.com"
export GITLAB_REPOS="project1/repo1,project2/repo1"
export SLACK_TOKEN="<super_secret_slack_token>"
export SLACK_CHANNEL="my_slack_room"
export USER_WHITELIST="user1,user2"
```

GitLab configuration is optional.

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


