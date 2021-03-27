# inti-activity
This tool, written in Go, constantly monitors Activity feed on intigriti.com Dashboard page and sends Slack/Discord notifications on new activities.

As Intigriti does not provide official API for reasearchers this tool mimics the login process to connect to API.
That's why you need to provide your full Intigriti login credentials at start.

**NOTE:** It doesn't work yet if you have 2FA authentication enabled.
## Installation
Install it from GitHub repository:
```
go get github.com/0xJeti/intitools/cmd/inti-activity
```

### Get Webhook URL
You need to create a Webhook to your Discord / Slack channel:
  * Slack - see [Sending messages using Incoming Webhooks](https://api.slack.com/messaging/webhooks)
  * Discord - See [Intro to Webhooks](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
  
## Run
`inti-activity` should run as a background process. Use your favourite way of daemonizing the process (`nohup`, `screen`, `tmux`, `systemd` etc.)

Usage of `inti-activity`:
```
  -config:      Path to config file
  -username:    Intigriti username (e-mail)
  -password:    Intigriti password
  -webhook:     Webhook URL
  -type:        Webhook type [slack|discord]
  -tick:        Ticking interval (dafault 60s)
  -last:        Number of activity entries sent on start (for debugging)
```

Alternatively you can provide config file with all parameters defined: (name it e.g. *monitor.conf*):

```
username YOUR_EMAIL
password YOUR_PASSWORD
webhook WEBHOOK
type    slack
tick 60s
```

and run the monitor with `-config` parameter:
```
inti-activity -config monitor.conf
```


