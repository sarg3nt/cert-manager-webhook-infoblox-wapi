# cspell:ignore KUBEBUILDER
GO ?= $(shell which go)
OS ?= $(shell $(GO) env GOOS)
ARCH ?= $(shell $(GO) env GOARCH)

IMAGE_NAME := ghcr.io/sarg3nt/cert-manager-webhook-infoblox-wapi
IMAGE_TAG := 2.0.0-beta1
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD | sed 's/[\/_]/-/g')

OUT := $(shell pwd)/_out

# K8s version for envtest (can be overridden, e.g., make test ENVTEST_K8S_VERSION=1.31.x)
ENVTEST_K8S_VERSION ?= 1.31.x

HELM_FILES := $(shell find charts/cert-manager-webhook-infoblox-wapi)

# Install setup-envtest tool
SETUP_ENVTEST = $(shell pwd)/bin/setup-envtest
.PHONY: setup-envtest
setup-envtest: $(SETUP_ENVTEST)
$(SETUP_ENVTEST): | bin
	GOBIN=$(shell pwd)/bin $(GO) install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: test
test: setup-envtest
	@KUBEBUILDER_ASSETS_PATH="$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" && \
	TEST_ASSET_ETCD="$$KUBEBUILDER_ASSETS_PATH/etcd" \
	TEST_ASSET_KUBE_APISERVER="$$KUBEBUILDER_ASSETS_PATH/kube-apiserver" \
	TEST_ASSET_KUBECTL="$$KUBEBUILDER_ASSETS_PATH/kubectl" \
	$(GO) test -v .

.PHONY: test-coverage
test-coverage: setup-envtest
	@KUBEBUILDER_ASSETS_PATH="$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" && \
	TEST_ASSET_ETCD="$$KUBEBUILDER_ASSETS_PATH/etcd" \
	TEST_ASSET_KUBE_APISERVER="$$KUBEBUILDER_ASSETS_PATH/kube-apiserver" \
	TEST_ASSET_KUBECTL="$$KUBEBUILDER_ASSETS_PATH/kubectl" \
	$(GO) test -coverprofile=coverage.out -covermode=atomic .

.PHONY: test-race
test-race: setup-envtest
	@KUBEBUILDER_ASSETS_PATH="$(shell $(SETUP_ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" && \
	TEST_ASSET_ETCD="$$KUBEBUILDER_ASSETS_PATH/etcd" \
	TEST_ASSET_KUBE_APISERVER="$$KUBEBUILDER_ASSETS_PATH/kube-apiserver" \
	TEST_ASSET_KUBECTL="$$KUBEBUILDER_ASSETS_PATH/kubectl" \
	$(GO) test -race -v .

.PHONY: clean
clean:
	rm -rf _test $(OUT) bin coverage.out

.PHONY: build
build:
	go mod tidy
	CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

.PHONY: build-container
build-container:
	go mod tidy
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)-$(GIT_BRANCH)" .

.PHONY: push-container
push-container: 
	docker push "$(IMAGE_NAME):$(IMAGE_TAG)-$(GIT_BRANCH)"

.PHONY: helm
helm: $(OUT)/rendered-manifest.yaml

$(OUT)/rendered-manifest.yaml: $(HELM_FILES) | $(OUT)
	helm template \
		cert-manager-webhook-infoblox-wapi \
		--set image.repository=$(IMAGE_NAME) \
		--set image.tag=$(IMAGE_TAG) \
		charts/cert-manager-webhook-infoblox-wapi > $@

_test $(OUT) bin:
	mkdir -p $@
