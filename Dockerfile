# VER=1.11.0 && IMAGE="ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi" && docker build . -t ${IMAGE}:${VER} && docker push ${IMAGE}:${VER} && docker tag ${IMAGE}:${VER} ${IMAGE}:latest && docker push ${IMAGE}:latest

# https://hub.docker.com/_/golang/
FROM golang:1.25.1 AS build_deps

LABEL org.opencontainers.image.source=https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_deps AS build

COPY . .
RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/sarg3nt/cert-manager-webhook-infoblox-wapi"
LABEL org.opencontainers.image.description="cert-manager Infoblox WAPI Webhook"
LABEL org.opencontainers.image.licenses="Apache-2.0"
# Removes upgrade artifacts to make the image smaller
ADD https://curl.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=build /workspace/webhook /webhook

ENTRYPOINT ["/webhook", "-v=4"]
