# Overview

A dummy subscriber that can be used for debugging and troubleshooting event delivery

It exposes a **REST API** `POST /v1/events` for receiving the events.

## Deploy Subscriber

```bash
kubectl apply -f ./config/subscriber-deployment.yaml
```

## Knative deployments

* Annotate your namespace

```bash
kubectl label namespace default knative-eventing-injection=enabled
```

* Deploy Knative broker

```bash
kubectl apply -f ./config/knative-broker.yaml
```

* Deploy Knative trigger

```bash
kubectl apply -f ./config/knative-trigger.yaml
```

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
