# Overview

This project contains code and configurations to explore various eventing frameworks such as

* Kyma eventing
* Knative eventing
* more to come

General idea is to:
 
* Deploy a susbcriber app that can consume the events
* Do the necessary configuration required for the framework.
* Send an event via a K8S pod using curl.

## Deploy Subscriber

```bash
kubectl apply -f ./config/subscriber-deployment.yaml
```

## Framework configuration

* [knative-eventing](./knative-eventing.md)


## Publish event

* Start a publisher pod with curl

```bash
kubectl run -i --tty --image appropriate/curl sample-publisher --restart=Never --rm /bin/sh
```

* Publish an event

```bash
curl -v "http://sample-broker-broker.default.svc.cluster.local/" \
  -X POST \
  -H "X-B3-Flags: 1" \
  -H "CE-CloudEventsVersion: 0.1" \
  -H "CE-EventType: order.created" \
  -H "CE-EventTime: 2018-04-05T03:56:24Z" \
  -H "CE-EventID: 45a8b444-3213-4758-be3f-540bf93f85ff" \
  -H "CE-Source: sample-external-solution" \
  -H 'Content-Type: application/json' \
  -d '{ "much": "wow" }'
```
## Verify the subscriber logs

```bash
kubectl logs -l app=dummy-subscriber -c dummy-subscriber
```