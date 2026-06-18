# Use make to manually build the container.
# https://hub.docker.com/_/golang/
FROM golang:1.26.4@sha256:792443b89f65105abba56b9bd5e97f680a80074ac62fc844a584212f8c8102c3 AS build_deps

LABEL org.opencontainers.image.source=https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi"
LABEL org.opencontainers.image.description="cert-manager Infoblox WAPI Webhook"
LABEL org.opencontainers.image.licenses="Apache-2.0"
# Removes upgrade artifacts to make the image smaller
ADD https://curl.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["/webhook", "-v=4"]
