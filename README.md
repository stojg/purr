# purr

Send open Github requests to a Slack channel

## installation

`go get -u github.com/stojg/purr`

## configuration

```
export GITHUB_TOKEN="<your personal token here"
export GITHUB_REPOS="user_org/repo1,user_org/repo2"
export SLACK_TOKEN="super_siikrit_slack_token"
export SLACK_CHANNEL="my_slack_room"
```

## run it

`purr`

This is a one shot action, so you might want to put into a cron or a [systemd timer unit](https://wiki.archlinux.org/index.php/Systemd/Timers)

## End result:

![example.png](./_docs/example.png)
