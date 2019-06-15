# Kwatchman

Watch for k8s workload resources and notify changes on Instant Messaging services

> :warning: The production go code is still a working in progress

## Build
```
$ VERSION=1.0.0 make build
```

### Docker image
```
$ VERSION=1.0.0 make docker-image
```
## Configure slack

In order to post messages with kwaychman to slack in a channel you have to 

1. [Create an slack application](https://api.slack.com/apps/new), you can call it `kwatchamn`
2. Create an Incoming Webhook, the url will be use to configure kwatchman later on

Both steps are pretty much the same as if you follow [slack's hello wolrd tutorial](https://api.slack.com/tutorials/slack-apps-hello-world)



