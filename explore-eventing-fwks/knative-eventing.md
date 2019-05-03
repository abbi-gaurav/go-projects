## Configure knative eventing

**Prerequisite**

- knative serving has been installed

## Steps

* Install Knative eventing

```bash
kubectl apply -f ./config/kn-eventing-0.5.yaml
```

* Create in-memory CCP

```bash
kubectl apply -f ./config/in-memory-channel-0.5.yaml
```

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