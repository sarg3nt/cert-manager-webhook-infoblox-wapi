# Use make to manually build the container.
# https://hub.docker.com/_/golang/
#
# The build stage is pinned to the native BUILDPLATFORM and cross-compiles to the
# requested TARGETOS/TARGETARCH. Because the webhook is a static CGO_ENABLED=0
# binary, cross-compilation produces identical output to a native build while
# running the Go toolchain at full native speed. This avoids emulating the whole
# Go build under QEMU for non-amd64 platforms (e.g. linux/arm64), which is
# extremely slow for a dependency tree this size.
FROM --platform=$BUILDPLATFORM golang:1.26.4@sha256:f96cc555eb8db430159a3aa6797cd5bae561945b7b0fe7d0e284c63a3b291609 AS build_deps

LABEL org.opencontainers.image.source=https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

# Provided automatically by buildx for the platform currently being built.
ARG TARGETOS
ARG TARGETARCH

COPY . .
# go.mod / go.sum are committed and verified tidy in CI, so do not run
# `go mod tidy` here -- it would re-resolve the module graph (and hit the
# network) on every image build. Cross-compile straight to the target platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi"
LABEL org.opencontainers.image.description="cert-manager Infoblox WAPI Webhook"
LABEL org.opencontainers.image.licenses="Apache-2.0"
# Removes upgrade artifacts to make the image smaller
ADD https://curl.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["/webhook", "-v=4"]
