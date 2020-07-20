FROM golang:1.14 as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .

RUN go mod download

COPY main.go .
COPY internal internal

RUN ls /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o sample-broker .

FROM scratch
WORKDIR /app
COPY --from=builder /app/sample-broker /app/

EXPOSE 8080
ENTRYPOINT ["/app/sample-broker"]