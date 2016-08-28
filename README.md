# purr

Send open Github requests to a Slack channel

## installation

`go get -u github.com/stojg/purr`

... or download one of the binaries for your platform from https://github.com/stojg/purr/releases

If you want to cross-compile for other platform, see http://golangcookbook.com/chapters/running/cross-compiling/

## configuration

```
export GITHUB_TOKEN="<super_secret_github token>"
export SLACK_TOKEN="<super_secret_slack_token>"
export GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
export SLACK_CHANNEL="my_slack_room"
```

## run it

`purr`

This is a one shot action, so you might want to put into a cron or a [systemd timer unit](https://wiki.archlinux.org/index.php/Systemd/Timers)

example `/etc/cron.d/purr` cron that runs purr 8am every day:

```
GITHUB_TOKEN="<super_secret_github token>"
SLACK_TOKEN="<super_secret_slack_token>"
GITHUB_REPOS="user_org/repo1,user_org/repo2" # comma separated
SLACK_CHANNEL="my_slack_room"
0 8 * * * root /usr/bin/purr 
```

## End result:

![example.png](./_docs/example.png)
