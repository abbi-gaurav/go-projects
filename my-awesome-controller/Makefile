RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=github.com/abbi-gaurav/go-learning-projects/my-awesome-controller
APP?=my-awesome-controller
PORT?=8000
GOOS?=linux
GOARCH?=amd64
CONTAINER_IMAGE?=gabbi/go-learning-projects/abbi-gaurav/${APP}

clean:
	rm -rf ${APP}

vet:
	go list ./... | grep -v 'pkg/client' | xargs go vet

run-test:
	go list ./... | grep -v 'pkg/client' | xargs go test -v -count=1

compile:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

build: clean vet run-test compile

image: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

helm-clean:
	helm delete ${APP} --purge

helm-deploy:
	helm install ./charts/my-awesome-controller --name ${APP} --debug

helm-deploy-dry-run:
	helm install ./charts/my-awesome-controller --name ${APP} --debug --dry-run
