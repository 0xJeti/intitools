# intitools
Collection of tools to interact with https://app.intigriti.com website.


# inti-activity
Monitors Activity feed on Dashboard page and sends Slack/Discord notifications on new activities.
NOTE: It doesn't work yet with 2FA enabled.
## Installation
### Install from GitHub repository:
```
go get github.com/0xJeti/intitools/cmd/inti-activity
```

### Get Webhook URL
  * Slack - see [Sending messages using Incoming Webhooks](https://api.slack.com/messaging/webhooks)
  * Discord - See [Intro to Webhooks](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
  
### Create config file (e.g. *monitor.conf*)
```
username YOUR_EMAIL
password YOUR_PASSWORD
webhook WEBHOOK
type    slack
tick 60s
```
* `tick` - duration between checks
* `webhook` - Slack/Discord Webhook URLS
* `type` - Webhook type: `slack` | `discord`

## Run
```
inti-activity -config monitor.conf
```
