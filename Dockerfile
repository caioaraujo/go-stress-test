FROM golang:1.22.2-alpine AS builder
ENV CGO_ENABLED=1

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod verify && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o stress-test

ENTRYPOINT [ "./stress-test" ]