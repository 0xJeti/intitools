Small tool, written in Go, that constantly monitors Activity feed on https://app.intigriti.com Dashboard page and sends Slack/Discord notifications on new activities. 

### Discord notifications
![Discord notification](https://github.com/0xJeti/intitools/raw/main/image/discord-notify.png)
### Slack notifications
![Slack notification](https://github.com/0xJeti/intitools/raw/main/image/slack-notify.png)

It is capable of showing differences between changes (new domains added/removed, updated In Scope, Out of Scope, descriptions etc.):

![Discord changes](https://github.com/0xJeti/intitools/raw/main/image/discord-changes.png)


# Installation
> As Intigriti does not provide official API for researchers this tool mimics the login process to connect to API. That's why you need to provide your full Intigriti login credentials at start.

Install / update from GitHub repository:
```
go install github.com/0xJeti/intitools/cmd/inti-activity@latest
```

## Webhook URL
You need to create a Webhook to your Discord / Slack channel:
  * Slack - see [Sending messages using Incoming Webhooks](https://api.slack.com/messaging/webhooks)
  * Discord - See [Intro to Webhooks](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)

## 2FA secret
If you use Two-factor authentication you need to obtain your 2FA secret. If your password manager is not giving you this information (e.g. Google Authenticator) you need to disable 2FA and enable it again to get the secret:
![2FA Secret](https://github.com/0xJeti/intitools/raw/main/image/2fa-secret.jpg)

**WARNING:** Copy the secret, uppercase it and remove the spaces before use. 
# Usage
`inti-activity` should be executed as a background process. Use your favourite way of daemonizing the process (`nohup`, `screen`, `tmux`, `systemd` etc.)

Usage of `inti-activity`:
```
  -config:      Path to config file (optional)
  -username:    Intigriti username (e-mail)
  -password:    Intigriti password
  -secret:      Intigriti 2FA secret (optional) 
  -webhook:     Webhook URL
  -type:        Webhook type [slack|discord]
  -tick:        Ticking interval (optional, dafault 60s)
  -last:        Number of activity entries sent on start (optional, for debugging)
```

You can provide all mandatory parameters via command line arguments.
Alternatively you can create a config file with some parameters defined: (name it e.g. *monitor.conf*):

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
