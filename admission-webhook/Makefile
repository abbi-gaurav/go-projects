RELEASE?=0.0.1
PROJECT?=github.com/go-learning-projects/admission-webhook
APP?=admission-webhook
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

build: clean test
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE}" \
		-o ${APP}

image: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

helm-clean:
	helm delete ${APP} --purge

helm-deploy:
	helm install ./charts/${APP} --name ${APP} --debug

helm-deploy-dry-run:
	helm install ./charts/${APP} --name ${APP} --debug --dry-run
