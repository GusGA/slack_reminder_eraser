# SLACK Reminder eraser

Inspirated by this blog [post](https://i.usedtocode.com/2017/02/11/how-to-batch-remove-all-your-reminders-from-slack/), This go lang script deletes all reminders asociated to an user.

## How to use

`go get https://github.com/gusga/slack_reminder_eraser`

`cd $GOPATH/src/github.com/gusga/slack_reminder_eraser`

`env SLACK_TOKEN=token SLACK_USER=user go run main.go`

This script can be used as binary too (Duhh!!!)

`env GOOS=OS GOARCH=ARCH go build -o eraser -v github.com/gusga/slack_reminder_eraser`

`env SLACK_TOKEN=token SLACK_USER=user eraser`


LICENSE MIT
