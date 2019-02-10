"""Watch multiple K8s event streams without threads."""
import asyncio
import json
import copy
import difflib
import os, sys, traceback
import requests
from time import time
from abc import ABC, abstractmethod
from kubernetes_asyncio import client, config, watch


async def watch_deployments(queue):
    v1 = client.AppsV1Api()
    async for event in watch.Watch().stream(v1.list_deployment_for_all_namespaces):
        await queue.put(event)

async def watch_daemonsets(queue):
    v1 = client.AppsV1Api()
    async for event in watch.Watch().stream(v1.list_daemon_set_for_all_namespaces):
        await queue.put(event)

async def watch_statefulsets(queue):
    v1 = client.AppsV1Api()
    async for event in watch.Watch().stream(v1.list_stateful_set_for_all_namespaces):
        await queue.put(event)

async def watch_cronjobs(queue):
    v1 = client.BatchV1beta1Api()
    async for event in watch.Watch().stream(v1.list_cron_job_for_all_namespaces):
        await queue.put(event)


async def sync_up(has_synced):
    # The naive approach is to just use a timer...
    start_time = time()
    while time() - start_time < 5:
        await asyncio.sleep(0.5)

    # Simple types assigment is an atomic operation
    has_synced.done = True
    print('sync-up')

async def consume_events(queue, has_synced, func):
    storage = {}
    count = 0
    while True:
        event = await queue.get()
        count += 1
        func(event, storage, has_synced, count)

class HasSynced(object):
    def __init__(self):
        self.done = False

class WorkloadsResource(ABC):
    def __init__(self, event):
        self.obj       = event['object']
        self.kind      = event['object'].kind
        self.name      = event['object'].metadata.name
        self.namespace = event['object'].metadata.namespace
        self.version   = event['object'].metadata.resource_version

    @property
    def idx_obj(self):
        return '{}/{}/{}'.format(self.namespace, self.kind, self.name)

    @property
    def clean(self):
        # TODO: Should we start instead with only spec?
        meta = copy.copy(self.obj.metadata).to_dict()
        if 'annotations' in meta and meta['annotations'] != None:
            meta['annotations'] = None

        meta['resource_version'] = None
        meta['generation'] = None

        spec = str(self.obj.spec)
        return '{}\n{}'.format(str(meta), spec)

class Deployment(WorkloadsResource):
    pass

class DaemonSet(WorkloadsResource):
    pass

class StatefulSet(WorkloadsResource):
    pass

class CronJob(WorkloadsResource):
    pass

def _notify_event(action, kind, name, event_count, diff):
    print(action, kind, name, event_count)
    print(diff)

    webhook_url = os.environ['SLACK_WEBHOOK_URL']

    short_message = "{} {}/{}".format(action, kind, name)
    if action == 'ADDED': color = "#7CD197"
    if action == 'DELETED': color = "#ff0000"
    if action == 'MODIFIED': color = "#ff9900"

    slack_data = {
        "attachments": [
            {
                "fallback": short_message,
                "title": short_message,
                "text": '```{}```'.format(diff),
                "color": color
            }
        ]
    }

    response = requests.post(
        webhook_url, data=json.dumps(slack_data),
        headers={'Content-Type': 'application/json'}
    )
    if response.status_code != 200:
        print('ERROR: slack webhook status code {}'.format(response.status_code))


def _compare_manifests(a, b):
    diff = '\n'.join(
        filter(
            lambda x: not x.startswith('@@') and not x.startswith('---') and not x.startswith('+++'),
            difflib.unified_diff(
                str(a).split('\n'), 
                str(b).split('\n'),
                n=0, 
                lineterm=''
            )
        )
    )
    return diff

def _handle_event(event, storage, has_synced, count):
    try:
        kind      = event['object'].kind
        action    = event['type']

        klasses = {
            'Deployment': Deployment, 
            'DaemonSet': DaemonSet, 
            'StatefulSet': StatefulSet, 
            'CronJob': CronJob
        }

        if kind in klasses:
            obj = klasses[kind](event)

        else:
            print('WARNING: Unknown resource kind {}'.format(kind))
            return

        new_object = obj.clean
        diff = None

        if has_synced.done:
            if action == 'ADDED' or action == 'DELETED':
                diff = _compare_manifests('', new_object)

            if action == 'MODIFIED':
                if obj.idx_obj in storage:
                    stored_object = storage[obj.idx_obj]
                    if stored_object != new_object: # TODO: Could dicts come with a different order and make equality to fail?
                        diff = _compare_manifests(stored_object, new_object)

        if diff:
            _notify_event(action, kind, obj.name, count, diff)

        storage[obj.idx_obj] = new_object

    except Exception as e:
        print(e)
        traceback.print_exc(file=sys.stdout)


if __name__ == '__main__':

    loop = asyncio.get_event_loop()
    queue = asyncio.Queue(loop=loop)
    loop.run_until_complete(config.load_kube_config())

    has_synced = HasSynced()
    tasks = [
        asyncio.ensure_future(watch_deployments(queue)),
        asyncio.ensure_future(watch_statefulsets(queue)),
        asyncio.ensure_future(watch_daemonsets(queue)),
        asyncio.ensure_future(watch_cronjobs(queue)),
        asyncio.ensure_future(consume_events(queue, has_synced, _handle_event)),
        asyncio.ensure_future(sync_up(has_synced)),
    ]

    loop.run_until_complete(asyncio.wait(tasks))
    loop.close()
