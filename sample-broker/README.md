# Overview

## Setup environment

```shell script
export NS={namespace}
export BA_USER={basic auth user}
export BA_PASSWORD={basic auth password}
export CLUSTER_DOMAIN={cluster domain}
```

## Deploy Broker

```shell script
  make deploy-broker
```

## Test the broker

* Get the available services

```shell script
curl --request GET "https://sample-broker.${CLUSTER_DOMAIN}/v2/catalog" \
  --header 'X-Broker-API-Version: 2.14' \
  --user $BA_USER:$BA_PASSWORD
```

* Provision an example service

```shell script
curl --location --request PUT "https://sample-broker.${CLUSTER_DOMAIN}/v2/service_instances/test-1" \
  --header 'X-Broker-API-Version: 2.14' \
  --header 'instance_id: test-1' \
  --user $BA_USER:$BA_PASSWORD \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "service_id": "123-123-123-123",
      "plan_id": "default",
      "context": {},
      "parameters": {
          "serviceInstanceName": "my-order-api",
          "auth" : {
              "baUser" : "user1",
              "baPassword" : "secret"
          }
      }
  }'
```

* Clean up example

```shell script
make clean-up-example
```

### Delete broker

```shell script
make delete-broker
```