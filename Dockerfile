# syntax=docker/dockerfile:1

# VER=1.6.0 && IMAGE="ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi" && docker build . -t ${IMAGE}:${VER} && docker push ${IMAGE}:${VER} && docker tag ${IMAGE}:${VER} ${IMAGE}:latest && docker push ${IMAGE}:latest

# https://hub.docker.com/_/golang/
FROM golang:1.23-alpine3.20@sha256:6a8532e5441593becc88664617107ed567cb6862cb8b2d87eb33b7ee750f653c AS build_deps

LABEL org.opencontainers.image.source=https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi

RUN apk add --no-cache git

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .
RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

# https://hub.docker.com/_/alpine/
FROM alpine:3@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099 AS alpine-upgraded

# Update all apk packages and install ca-certificates
RUN apk upgrade --no-cache && \
  apk add --no-cache ca-certificates

# Main image
FROM scratch
# Removes upgrade artifacts to make the image smaller
COPY --from=alpine /etc/ssl/certs /etc/ssl/certs
# Copy over the compiled webhook executable from the builder.
COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["./webhook", "-v=4"]
