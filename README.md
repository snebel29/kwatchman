# Kwatchman

Watch for k8s workload resources and notify changes on Instant Messaging services

> :warning: The production go code is still in alpha version

## Install

### Using helm chart

1. [Install helm](https://helm.sh/docs/using_helm/)
2. Clone this repository
```
$ git clone https://github.com/snebel29/kwatchman.git
```
3. Install the chart, for now this have to be installed from your local file system, but will be published into https://github.com/helm/charts soon
```
$ cd kwatchman/build/chart/kwatchman
$ helm install -n kwatchman .
```

Docker images available can be found in [snebel29/kwatchman dockerhub](https://cloud.docker.com/repository/docker/snebel29/kwatchman)

### Configure slack

In order to post messages with kwaychman to slack in a channel you have to 

1. [Create an slack application](https://api.slack.com/apps/new), you can call it `kwatchamn`
2. Create an Incoming Webhook, the url will be use to configure kwatchman later on

Both steps are pretty much the same as if you follow [slack's hello wolrd tutorial](https://api.slack.com/tutorials/slack-apps-hello-world)

## Usage
```
usage: kwatchman --slack-webhook=SLACK-WEBHOOK [<flags>]

Flags:
  -h, --help          Show context-sensitive help (also try --help-long and --help-man).
  -c, --cluster-name="undefined"  
                      Name of k8s cluster where kwatchman is running, use for notification purposes only
  -n, --namespace=""  k8s namespace where to get resources from: default to all
  -k, --kubeconfig="/Users/svennebel/.kube/config"  
                      kubeconfig path for running out of k8s
  -w, --slack-webhook=SLACK-WEBHOOK  
                      The slack webhook url (Required)
      --version       Show application version.
```

## Development
### Build
```
$ VERSION=1.0.0 make build
```

#### Docker image
```
$ VERSION=1.0.0 make docker-image
```
