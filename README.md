# Kwatchman

Kwatchman is a tool that watch for resources events and manifest changes within kubernetes clusters and triggers chain of handlers with its event information.

For example, whenever the container image tag version of a deployment changes you trigger a notification on slack, or you simply log the event into your corporate structured logging platform for future troubleshooting and processing.

You can check the project roadmap within its [kwatchman repository project](https://github.com/snebel29/kwatchman/projects/1)

  ![](img/demo.gif)

## Install

Kwatchman is delivered in docker images, which can be found in [snebel29/kwatchman dockerhub](https://hub.docker.com/r/snebel29/kwatchman/tags) but you can also run it from your any computer with a valid kubeconfig file even with special authentication requirements such as cloud managed clusters solutions.

You can create your own kubernetes manifest and deploy althought the usual way of installing kwatchman should be using the kubernetes official package manager [helm](https://helm.sh/) 

### Install using helm chart

1. [Install helm](https://helm.sh/docs/using_helm/)
2. Install the chart, for now this have to be installed from your local file system, but will be published into https://github.com/helm/charts soon
```
$ git clone https://github.com/snebel29/kwatchman.git
$ cd kwatchman/build/chart/kwatchman
$ helm install -n kwatchman .
```

### Configuration

Kwatchman uses [viper](https://github.com/spf13/viper) to read from a config file in any of its accepted format, although since all the tests and examples were created using toml, this is the recommended format.

The format is pretty simple, with two main sections to configure resources and handlers, detailed information on the configuration can be found in the default [config.toml](./config.toml) file within the root of this repository.


#### Resources

Define the kubernetes resources to watch, not all resources are available to watch although they are continuosly added, create an [issue](https://github.com/snebel29/kwatchman/issues) or contribute yourself to get more resources added

> :information_source: Resources should handle apiGroup deprecation and removal transparently for the user while using last stable kwatchman versions

#### Handlers
Handlers is what makes kwatchman powerfull, takes as input all the related event information (action, k8s manifest, etc..) and execute some code using it, they can be used for notifiying to instant message services such as Slack or to simply log the events into your logging system, handlers can be chained passing the responsability to keep the execution to the next handler.

Currently only a hand of handlers can be used, but there is plans to allow to building your own through plugins, webhooks and local executor, match the interface, keep code with reasonable quality and consistency and you can as well pull request for adding useful general purpose handlers.

> :information_source: Upon a handler error, the whole chain is automatically retried up to 3 times

###### The diff handler
This is the most powerful handler, and should be for almost all use cases the base handler as it filters noisy events reporting only the changes to manifests.

###### The log handler
This is used often for testing but for recording reported changes right after the diff handler as well, enriching your logs and metrics with fresh high level events that can be used for root cause analysis either for humans or machines (AIOps)

###### The Slack handler
Have you ever being bitten by a change that a colleague never reported? running this handler after the `diff` handler you will get notifications into your slack channel about any change in your cluster. 

####### Configure slack

In order to post messages with kwatchman to slack in a channel you have to 

1. [Create an slack application](https://api.slack.com/apps/new), you can call it `kwatchamn`
2. Create an Incoming Webhook, the url will be use to configure kwatchman later on

Both steps are pretty much the same as if you follow [slack's hello wolrd tutorial](https://api.slack.com/tutorials/slack-apps-hello-world)

## Compatibility matrix
Kwatchman uses [go-client](https://github.com/kubernetes/client-go) and a forked version of [kooper](https://github.com/snebel29/kooper) and it's therefore coupled to their version compatibility, new released may be required in order to fully work with future versions, please check the compatibility matrix which is provided, any reported issue will be fixed in a best effor basis.

kwatchman version | k8s version |
|:----------:|:-------------:|
| v1.0.0 |  +1.11 |

## Similar projects

- [kubewatch](https://github.com/bitnami-labs/kubewatch)
- [chowkidar](https://github.com/stakater/Chowkidar)

## Development
### Dependencies

- [make](https://www.gnu.org/software/make/)
- [docker](https://www.docker.com/)
- [go](https://golang.org/dl/) >= 1.12 official builds currently uses 1.12
- [dep](https://github.com/golang/dep) will soon be replaced by go modules

### Testing
```
$ make test
```

### Build
```
$ VERSION=1.0.0 make build
```

#### Docker image
```
$ VERSION=1.0.0 make docker-image
```

### Release
Relase of new docker images is achieved by creating a git tag, and through circleci
```
$ git tag v1.0.0
$ git push origin master --tags
```

