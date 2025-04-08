# VER=1.6.0 && IMAGE="ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi" && docker build . -t ${IMAGE}:${VER} && docker push ${IMAGE}:${VER} && docker tag ${IMAGE}:${VER} ${IMAGE}:latest && docker push ${IMAGE}:latest

# https://hub.docker.com/_/golang/
FROM golang:1.24.2-alpine3.21@sha256:7772cb5322baa875edd74705556d08f0eeca7b9c4b5367754ce3f2f00041ccee AS build_deps

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
FROM alpine:3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c AS alpine-upgraded

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
