# intitools
Collection of tools to interact with https://app.intigriti.com website.


# inti-activity
Monitors Activity feed on Dashboard page and sends Slack notifications on new activities.
NOTE: It doesn't work yet with 2FA enabled.
## Installation
```
go get github.com/0xJeti/intitools/cmd/inti-activity
```

Create config file (e.g. *monitor.conf*)
```
username YOUR_EMAIL
password YOUR_PASSWORD
webhook SLACK_WEBHOOK
tick 60s
```
* `tick` - duration between checks
* `webhook` - Slack Webhook. See [Sending messages using Incoming Webhooks](https://api.slack.com/messaging/webhooks) for more details.
## Run
```
inti-activity -config monitor.conf
```
