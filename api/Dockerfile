FROM golang:1.22 as build

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY cmd ./cmd

RUN go build -o /app/main ./cmd/main

FROM alpine:latest

RUN apk add --no-cache redis

COPY --from=build /app/main /app/main

WORKDIR /app

CMD ["./main"]
