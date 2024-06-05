FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext

COPY ["go.mod", "go.sum", "./"]

RUN go mod download
