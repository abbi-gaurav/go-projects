RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=github.com/go-learning-projects/k8s-ready-service
APP?=advent
PORT?=8000
GOOS?=linux
GOARCH?=amd64
CONTAINER_IMAGE?=gabbi/go-learning/${APP}

clean:
	rm -rf ${APP}

vet:
	go tool vet .
test:
	go test -v -race ./...

build: clean vet test
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

image: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

helm-clean:
	helm delete ${APP} --purge

helm-deploy:
	helm install ./charts/k8s-ready-servcie --name ${APP} --debug

helm-deploy-dry-run:
	helm install ./charts/k8s-ready-servcie --name ${APP} --debug --dry-run
