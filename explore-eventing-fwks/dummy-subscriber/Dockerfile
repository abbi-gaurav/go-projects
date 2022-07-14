FROM golang:1.18.3 as builder

ENV GO111MODULE=on

WORKDIR /go/src/dummy-subscriber
COPY go.mod .

COPY main.go .
COPY internal internal

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o dummy-subscriber .

FROM scratch
WORKDIR /root/

COPY --from=builder /go/src/dummy-subscriber/dummy-subscriber .

EXPOSE 8080
ENTRYPOINT ["/root/dummy-subscriber"]
