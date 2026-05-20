# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY monolito ./monolito

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/gotodolist-monolito ./monolito

FROM alpine:3.22.4 AS runtime

RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=build /out/gotodolist-monolito /app/gotodolist-monolito

EXPOSE 8080

USER app

ENTRYPOINT ["/app/gotodolist-monolito"]
