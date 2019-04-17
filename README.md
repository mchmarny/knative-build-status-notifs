# knative-build-status-notifs

> The readme is still work in progress, when in doubt, check the `Makefile` for now

Simple Cloud Build status notification based on Cloud PubSub event source for Knative.

* [Cloud Build](https://cloud.google.com/cloud-build/)
* [Knative](https://github.com/knative/docs)
* [Knative Eventing](https://github.com/knative/docs/tree/master/eventing)
* [Demo App with Cloud Build Trigger](github.com/mchmarny/knative-gitops-using-cloud-build)
* This Cloud Build mobile notification using Pushover


## Setup

This setup assumes you already have

* GCP project
* Configured gcloud CLI
* Knative cluster with Eventing and PubSub event source configured
* Pushover account with app and user token

### Setup Notification Service

First, create a secret with the pushover tokens

```shell
kubectl create secret generic build-notif-secrets -n demo \
	--from-literal=SLACK_API_TOKEN=$SLACK_API_TOKEN
```

Build this notification service

```shell
gcloud builds submit \
	--project $GCP_PROJECT \
	--tag gcr.io/${GCP_PROJECT}/build-notif
```

And edit `config/service.yaml` with your GCP project ID

```yaml
container:
    image: gcr.io/YOUR_PROJECT_ID_HERE/build-notif:latest
    imagePullPolicy: Always
```

After your done all these steps, you should be able to see the `build-notif` status as `Running`

```shell
$: kubectl get pods -n demo

NAME                                                            READY     STATUS      RESTARTS   AGE
build-notif-00001-deployment-745db8fcb-xf97f                    3/3       Running     0          1h
```

### Setup Build Status Event Source

Assuming you have Knative Eventing and PubSub event source already configured, you can simply connect your PubSub queue that was created automatically by Cloud Build with this simple `GcpPubSubSource` event source manifest:

```yaml
apiVersion: sources.eventing.knative.dev/v1alpha1
kind: GcpPubSubSource
metadata:
  name: cloud-build-status-source
spec:
  googleCloudProject: s9-demo
  topic: cloud-builds
  gcpCredsSecret:
    name: google-cloud-key
    key: key.json
  sink:
    apiVersion: eventing.knative.dev/v1alpha1
    kind: Broker
    name: default
```

### Setup Event Trigger 

The above event source will publish events into the default broker in that namespace so all we have to do next is to create a `Trigger` that connects to the `Service` using this manifest: 

```yaml
apiVersion: eventing.knative.dev/v1alpha1
kind: Trigger
metadata:
  name: slacker-build-status-notifier
spec:
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1alpha1
      kind: Service
      name: build-notif
```


## Demo

To demo the notifications, simply create a release tag on the repo and you should receive a Pushover notification.

## Debug

1. Check Notification service logs

```shell
kail -d build-notif-00006-deployment -n demo -c user-container --since 2h
```


2. Check PubSub source logs

```shell
kail -p gcppubsub-cloud-build-status-source-... -n demo -c user-container --since 2h
```

> You can find the names of the pods below using

```shell
kubectl get pods -n demo
```

3. Check build status in Cloud Build

Notifications are triggered only when build actually runs so make sure it was triggered when you tagged your release

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
