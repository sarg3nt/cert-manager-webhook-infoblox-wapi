# syntax=docker/dockerfile:1

# VER=1.6.0 && IMAGE="ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi" && docker build . -t ${IMAGE}:${VER} && docker push ${IMAGE}:${VER} && docker tag ${IMAGE}:${VER} ${IMAGE}:latest && docker push ${IMAGE}:latest

# https://hub.docker.com/_/golang/
FROM golang:1.23-alpine3.20@sha256:c694a4d291a13a9f9d94933395673494fc2cc9d4777b85df3a7e70b3492d3574 AS build_deps

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
FROM alpine:3@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45 AS alpine-upgraded

# Update all apk packages and install ca-certificates
RUN apk upgrade --no-cache && \
  apk add --no-cache ca-certificates

# Main image
FROM scratch
# Removes upgrade artifacts to make the image smaller
COPY --from=alpine-upgraded / /
# Copy over the compiled webhook executable from the builder.
COPY --from=build /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook", "-v=4"]
