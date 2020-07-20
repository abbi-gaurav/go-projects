NS?=default
RELEASE?=0.0.1
APP?=sample-broker
GOARCH?=amd64
CONTAINER_IMAGE?=gabbi/${APP}:${RELEASE}
NAMESPACE?=dev
USERNAME?=admin
PASSWORD?=kyma4ever

clean:
	rm -rf ${APP}

vet:
	GO111MODULE=on go vet .

test:
	GO111MODULE=on go test -v -race ./...

build: clean vet test
	GO111MODULE=on CGO_ENABLED=0 GOARCH=${GOARCH} go build \
		-ldflags "-s -w" \
		-o bin/${APP}

build-image: build
	docker build -t $(CONTAINER_IMAGE) .

push-image:
	docker push $(CONTAINER_IMAGE)

tidy:
	GO111MODULE=on go mod tidy

init:
	GO111MODULE=on go mod init

set-rbac:
	kubectl -n ${NS} create serviceaccount sample-broker
	kubectl -n ${NS} create role sample-broker --verb=get,list,watch,create,delete,update --resource=deployments,services,apirules
	kubectl -n ${NS} create rolebinding sample-broker --role=sample-broker --serviceaccount=${NS}:sample-broker

set-auth:
	kubectl -n ${NS} create secret generic sample-broker --from-literal=USERNAME=${USERNAME} --from-literal=PASSWORD=${PASSWORD}

set-catalog:
	kubectl -n ${NS} create configmap broker-catalog --from-file=examples/catalog.yaml

pre-deploy: set-rbac set-auth set-catalog

deploy:
	kubectl -n ${NS} apply -f k8s/deployment.yaml

run-locally:
	docker run -p 8080:8080 \
		-v ${PWD}/examples:/etc/sample-broker \
		-v /Users/d066419/work/service-factory/evangelist/skrs:/etc/kubernetes \
		-e USER_NAME=${USER} -e PASSWORD=${PASSWORD} \
		-e KUBECONFIG=/etc/kubernetes/marco-demo-canary.yaml \
		gabbi/sample-broker:0.0.1
