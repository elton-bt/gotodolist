# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY desacoplado/backend ./desacoplado/backend

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG APP_VERSION=dev

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X github.com/elton-bt/gotodolist/internal/config.version=${APP_VERSION}" -o /out/gotodolist-api ./desacoplado/backend

FROM alpine:3.22.4 AS runtime

RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=build /out/gotodolist-api /app/gotodolist-api

EXPOSE 8081

USER app

ENTRYPOINT ["/app/gotodolist-api"]
