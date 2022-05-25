APP?=ultimate-go-programming
GOARCH?=amd64

clean:
	rm -rf ${APP}

vet:
	GO111MODULE=on go vet ./...

test:
	GO111MODULE=on go test -v -race ./...

build: clean
	GO111MODULE=on CGO_ENABLED=0 GOARCH=${GOARCH} go build \
		-ldflags "-s -w" \
		-o ${APP}

tidy:
	GO111MODULE=on go mod tidy