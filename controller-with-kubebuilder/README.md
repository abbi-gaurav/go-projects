# Description

A sample controller built using kubebuilder. 

The controller is doing the CRUD operations to a dummy in-memory db based on the custom resource.

## Run the app against minikube

```bash
# install CRDs into cluster
make install
make run
kubectl apply -f ./config/samples/
```

## Deploy against minikube

```bash
# install CRDs into cluster
eval $(minikube docker-env)
make docker-build
make install
make deploy
kubectl apply -f ./config/samples/
```