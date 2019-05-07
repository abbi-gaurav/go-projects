# Overview

This guide references the steps to demonstrate the changing the default channel configuration. It starts with a state where no-default is specified (This is similar to how Kyma is currently deployed).

**Prerequisites**

- Knative serving is installed

## Steps

* Deploy knative eventing

    ```bash
    kubectl apply -f ./config/kn-eventing-0.5.yaml
    ```

* Deploy in memory CCP

    ```bash
    kubectl apply -f ./config/in-memory-ccp-0.5.yaml
    ```

* Update the `default-channel-webhook` to use in-memory as default

    ```bash
    kubectl apply -f ./config/in-memory-as-default-config.yaml
    ```

* create a channel with empty spec

    ```bash
    kubectl apply -f ./config/channel-with-in-memory-as-default.yaml
    ```

* check the channel yaml in cluster and verify that the spec contains `in-memory` as provisioner.

    ```bash
    kubectl get channels.eventing.knative.dev in-memory-as-default -o yaml
    ```

* Update the `default-channel-webhook` to use natss as default
   
    ```bash
    kubectl apply -f ./config/natss-as-default-config.yaml
    ```
   
* create a channel with empty spec
   
    ```bash
    kubectl apply -f ./config/channel-with-natss-as-default.yaml
    ```
   
* check the channel yaml in cluster and verify that the spec contains `natss` as provisioner.
   
    ```bash
    kubectl get channels.eventing.knative.dev natss-as-default -o yaml
    ```