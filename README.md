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
kubectl create secret generic slacknotifs-secrets -n demo \
	--from-literal=APP_TOKEN=$PUSHOVER_APP_TOKEN \
	--from-literal=USR_TOKEN=$PUSHOVER_USR_TOKEN
```

Build this notification service

```shell
gcloud builds submit \
	--project $GCP_PROJECT \
	--tag gcr.io/${GCP_PROJECT}/slacknotifs
```

And edit `config/service.yaml` with your GCP project ID

```yaml
container:
    image: gcr.io/YOUR_PROJECT_ID_HERE/slacknotifs:latest
    imagePullPolicy: Always
```

After your done all these steps, you should be able to see the `slacknotifs` status as `Running`

```shell
$: kubectl get pods -n demo

NAME                                                            READY     STATUS      RESTARTS   AGE
slacknotifs-00001-deployment-745db8fcb-xf97f                    3/3       Running     0          1h
```

### Setup Build Status Event Source

Assuming you have Knative Eventing and PubSub event source already configured, you can simply connect your PubSub queue that was created automatically by Cloud Build with the previously deployed service using `GcpPubSubSource` manifest

```yaml
apiVersion: sources.eventing.knative.dev/v1alpha1
kind: GcpPubSubSource
metadata:
  name: cloud-build-status-source
  namespace: demo
spec:
  googleCloudProject: s9-demo
  topic: cloud-builds
  gcpCredsSecret:
    name: google-cloud-key
    key: key.json
  sink:
    apiVersion: serving.knative.dev/v1alpha1
    kind: Service
    name: slacknotifs
    namespace: demo
```

## Demo

To demo the notifications, simply create a release tag on the repo and you should receive a Pushover notification.

## Debug

1. Check Notification service logs

```shell
kail -d slacknotifs-00006-deployment -n demo -c user-container --since 2h
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
