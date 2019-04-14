FROM golang:1.12 as builder

WORKDIR /go/src/dummy-subscriber
COPY main.go .
RUN ls ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o dummy-subscriber .

FROM scratch
WORKDIR /root/

COPY --from=builder /go/src/dummy-subscriber/dummy-subscriber .

EXPOSE 8080
ENTRYPOINT ["/root/dummy-subscriber"]
