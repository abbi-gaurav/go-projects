# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/
COPY db/     db/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/cmd/manager

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /go/src/github.com/abbi-gaurav/go-learning-projects/controller-with-kubebuilder/manager .
ENTRYPOINT ["/manager"]
