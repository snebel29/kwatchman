# Kwatchman prototype

This is a prototype built in python `3.7.1` and using `asyncio` to passively watch a kubernetes cluster and notify on slack for changes to resources used in workloads (Deployment, StatefulSet, DaemonSet and CronJob)

## Installation

```
$ git clone git@github.com:snebel29/kwatchman.git
$ cd kwatchman/prototype/
$ pyenv install -s 3.7.1
$ python -m venv .
$ source bin/activate
$ pip install -r requirements.txt
```

## Run
First of all you have to set your SLACK_WEBHOOK

```
$ export SLACK_WEBHOOK_URL='The webhook url'
```

> :warning: Remember to set appropiate cluster 

```
$ kubectl config use-context gke_snebel29-playground_us-central1-a_standard-cluster-1
```

Then simply run 
```
$ python async_kwatchman.py
```

and wait for `sync-up` message (should take 5 seconds) and the watcher is ready, watch your slack channel and you can edit any existing Deployment or deploy your own and see how notifications flow... 

Be aware that slack messages format is not polished yet and that this code may contain bugs.

## Ideas and Improvements

- Add k8s/service resource
- Improve has_synced function with a less naive approach (Is this necessary?)
- Enhance slack notification by providing cluster name and further context
- Enhance slack notification by using extended formating options such as upload file, etc.. (may require to use slack api bot)
- Enhance slack notification by uploading diff to GCS and provide link to it
- Enhance slack notification by watching for `kwachman` annotations with links to changes in external tooling (Gerrit, github) so Github PRs and/or Gerrit CL could be appended within CI and linked automatically into the nitification 
