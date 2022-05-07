# Copyright Yinan Li <cndoit18@outlook.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
REGISTRY := ghcr.io/cndoit18

GIT_COMMIT ?= $(shell git rev-parse HEAD)
APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
IMAGE_TAGS := $(APP_VERSION:v%=%)
IMAGE_NAME := lxcfs-manager
AGENT_IMAGE_NAME := lxcfs-agent
BUILD_TAG := dev

MANIFESTS_DIR ?= config
RBAC_DIR ?= $(MANIFESTS_DIR)/rbac
BOILERPLATE_FILE ?= ./hack/boilerplate.go.txt

GEN_CRD_OPTIONS ?= crd:trivialVersions=true
GEN_RBAC_OPTIONS ?= rbac:roleName=manager-role
GEN_WEBHOOK_OPTIONS ?= webhook
GEN_OBJECT_OPTIONS ?= object:headerFile=$(BOILERPLATE_FILE)
GEN_OUTPUTS_OPTIONS ?= output:rbac:artifacts:config=$(RBAC_DIR)

REPO = $(shell go list -m)
GOLDFLAGS="-X '$(REPO)/version.AppVersion=$(APP_VERSION)' -X '$(REPO)/version.GitCommit=$(GIT_COMMIT)'"

manifests: controller-gen
	$(CONTROLLER_GEN) paths="./..." $(GEN_RBAC_OPTIONS) $(GEN_WEBHOOK_OPTIONS) $(GEN_OBJECT_OPTIONS) $(GEN_OUTPUTS)

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	$(GOLANGCI_LINT) run --timeout 2m0s ./...

# Generate code
generate: manifests helm-generate ## Generate code, charts...

.PHONY: docker-build
docker-build: ## Build image
	docker build . -f Dockerfile --build-arg GOLDFLAGS=${GOLDFLAGS} -t $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG)
	docker build . -f Dockerfile.agent -t $(REGISTRY)/$(AGENT_IMAGE_NAME):$(BUILD_TAG)
	set -e; \
		for tag in $(IMAGE_TAGS); do \
			docker tag $(REGISTRY)/$(IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
			docker tag $(REGISTRY)/$(AGENT_IMAGE_NAME):$(BUILD_TAG) $(REGISTRY)/$(AGENT_IMAGE_NAME):$${tag}; \
	done

.PHONY: docker-push
docker-push: ## Push image
	set -e; \
		for tag in $(IMAGE_TAGS); do \
			docker push $(REGISTRY)/$(IMAGE_NAME):$${tag}; \
			docker push $(REGISTRY)/$(AGENT_IMAGE_NAME):$${tag}; \
	done

.PHONY: helm-push
helm-push: cr
	$(CR) index --config .github/.cr.yaml  --push

BIN ?= $(CURDIR)/.bin

$(BIN):
	mkdir -p "$(BIN)"

.PHONY: dev-tools
dev-tools: \
	controller-gen \
	golangci-lint \
	kustomize \
	helm-docs \
	cr

# find or download controller-gen
# download controller-gen if necessary
CONTROLLER_GEN_VERSION := 0.8.0
CONTROLLER_GEN := $(BIN)/controller-gen

.PHONY: controller-gen
controller-gen:
	@$(CONTROLLER_GEN) --version 2>&1 \
		| grep 'v$(CONTROLLER_GEN_VERSION)' \
	|| rm -f $(CONTROLLER_GEN)
	@$(MAKE) $(CONTROLLER_GEN)

$(CONTROLLER_GEN):
	$(MAKE) $(BIN)
	# https://github.com/kubernetes-sigs/controller-tools/tree/master/cmd/controller-gen
	go install 'sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_GEN_VERSION)'
	go build -mod=readonly -o $(CONTROLLER_GEN) sigs.k8s.io/controller-tools/cmd/controller-gen
	go mod tidy

# find or download helm-docs
# download helm-docs if necessary
HELM_DOCS_VERSION := 1.4.0
HELM_DOCS := $(BIN)/helm-docs

.PHONY: helm-docs
helm-docs:
	@$(HELM_DOCS) --version 2>&1 \
		| grep '$(HELM_DOCS_VERSION)' > /dev/null \
	|| rm -f $(HELM_DOCS)
	@$(MAKE) $(HELM_DOCS)

$(HELM_DOCS):
	$(MAKE) $(BIN)
	# https://github.com/norwoodj/helm-docs/tree/master#installation
	curl -sL "https://github.com/norwoodj/helm-docs/releases/download/v$(HELM_DOCS_VERSION)/helm-docs_$(HELM_DOCS_VERSION)_$$(uname -s)_x86_64.tar.gz" \
		| tar -xzC '$(BIN)' helm-docs

# find or download golangci-lint
# download golangci-lint if necessary
GOLANGCI_LINT := $(BIN)/golangci-lint
GOLANGCI_LINT_VERSION := 1.33.0

.PHONY: golangci-lint
golangci-lint:
	@$(GOLANGCI_LINT) version --format short 2>&1 \
		| grep '$(GOLANGCI_LINT_VERSION)' > /dev/null \
	|| rm -f $(GOLANGCI_LINT)
	@$(MAKE) $(GOLANGCI_LINT)

$(GOLANGCI_LINT):
	$(MAKE) $(BIN)
	# https://golangci-lint.run/usage/install/#linux-and-windows
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(BIN) 'v$(GOLANGCI_LINT_VERSION)'

# find helm or raise an error
.PHONY: helm
helm:
ifeq (, $(shell which helm 2> /dev/null))
	$(error Helm not found. Please install it: https://helm.sh/docs/intro/install/#from-script)
HELM=helm-not-found
else
HELM=$(shell which helm 2> /dev/null)
endif

# find or download kustomize
# download kustomize if necessary
KUSTOMIZE_VERSION := 3.8.7
KUSTOMIZE := $(BIN)/kustomize

.PHONY: kustomize
kustomize:
	@$(KUSTOMIZE) version --short 2>&1 \
		| grep 'kustomize/v$(KUSTOMIZE_VERSION)' > /dev/null \
	|| rm -f $(KUSTOMIZE)
	@$(MAKE) $(KUSTOMIZE)

$(KUSTOMIZE):
	$(MAKE) $(BIN)
	# https://kubectl.docs.kubernetes.io/installation/kustomize/binaries/
	curl -sSL "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v$(KUSTOMIZE_VERSION)/kustomize_v$(KUSTOMIZE_VERSION)_$$(go env GOOS)_$$(go env GOARCH).tar.gz" \
		| tar -xzC '$(BIN)' kustomize

# find or download kustomize
# download kustomize if necessary
CR_VERSION := 1.2.1
CR := $(BIN)/cr

.PHONY: cr
cr:
	@$(CR) version --short 2>&1 \
		| grep 'v$(CR_VERSION)' > /dev/null \
	|| rm -f $(CR)
	@$(MAKE) $(CR)

$(CR):
	$(MAKE) $(BIN)
	curl -sSL "https://github.com/helm/chart-releaser/releases/download/v$(CR_VERSION)/chart-releaser_$(CR_VERSION)_$$(go env GOOS)_$$(go env GOARCH).tar.gz" \
		| tar -xzC '$(BIN)' cr

####################
#    Helm chart    #
####################

RELEASE_VERSION ?= $(APP_VERSION:v%=%)
CHARTS_DIRECTORY      := charts
CHART_PROJECT_PATH := $(CHARTS_DIRECTORY)/lxcfs-on-kubernetes
CHART_TEMPLATE_PATH := $(CHART_PROJECT_PATH)/templates
DO_NOT_EDIT := Code generated by make. DO NOT EDIT.

$(CHART_TEMPLATE_PATH)/certificate.yaml: kustomize $(wildcard config/helm/certmanager/*) $(wildcard config/certmanager/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/certificate.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/certificate | \
	$(KUSTOMIZE) cfg grep --annotate=false 'kind=Certificate' | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/certificate.yaml

$(CHART_TEMPLATE_PATH)/issuer.yaml: kustomize $(wildcard config/helm/certmanager/*) $(wildcard config/certmanager/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/issuer.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/certificate | \
	$(KUSTOMIZE) cfg grep --annotate=false 'kind=Issuer' | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/issuer.yaml

$(CHART_TEMPLATE_PATH)/clusterrole.yaml: kustomize $(wildcard config/helm/rbac/*) $(wildcard config/rbac/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/clusterrole.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/rbac | \
	$(KUSTOMIZE) cfg grep --annotate=false "kind=ClusterRole" | \
	$(KUSTOMIZE) cfg grep --annotate=false --invert-match 'kind=ClusterRoleBinding' | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/clusterrole.yaml

$(CHART_TEMPLATE_PATH)/clusterrolebinding.yaml: kustomize $(wildcard config/helm/rbac/*) $(wildcard config/rbac/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/clusterrolebinding.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/rbac | \
	$(KUSTOMIZE) cfg grep --annotate=false "kind=ClusterRoleBinding" | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/clusterrolebinding.yaml

$(CHART_TEMPLATE_PATH)/mutatingwebhookconfiguration.yaml: kustomize $(wildcard config/helm/webhook/*) $(wildcard config/webhook/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/mutatingwebhookconfiguration.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/webhook | \
	$(KUSTOMIZE) cfg grep --annotate=false "kind=MutatingWebhookConfiguration" | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/mutatingwebhookconfiguration.yaml

$(CHART_TEMPLATE_PATH)/deployment.yaml: kustomize kustomize $(wildcard config/helm/manager/*) $(wildcard config/manager/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/deployment.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/manager | \
	$(KUSTOMIZE) cfg grep --annotate=false "kind=Deployment" | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/deployment.yaml

$(CHART_TEMPLATE_PATH)/daemonset.yaml: kustomize kustomize $(wildcard config/helm/manager/*) $(wildcard config/manager/*)
	echo '{{- /* $(DO_NOT_EDIT) */ -}}' > $(CHART_TEMPLATE_PATH)/daemonset.yaml
	echo '{{- if .Values.lxcfs.useDaemonset -}}' >> $(CHART_TEMPLATE_PATH)/daemonset.yaml
	$(KUSTOMIZE) build --reorder legacy config/helm/manager | \
	$(KUSTOMIZE) cfg grep --annotate=false "kind=DaemonSet" | \
	sed "s/'\({{[^}}]*}}\)'/\1/g" \
		>> $(CHART_TEMPLATE_PATH)/daemonset.yaml
	echo '{{- end -}}' >> $(CHART_TEMPLATE_PATH)/daemonset.yaml

$(CHART_PROJECT_PATH)/README.md: helm-docs $(CHART_PROJECT_PATH)/values.yaml $(CHART_PROJECT_PATH)/Chart.yaml
	$(HELM_DOCS) $(CHART_PROJECT_PATH)

define UPDATE_CHART_YAML
s#version: .*#version: $(RELEASE_VERSION)#
s#appVersion: .*#appVersion: $(RELEASE_VERSION)#
endef

define UPDATE_VALUES_YAML
s#lxcfs-manager:.*#lxcfs-manager:v$(IMAGE_TAGS)#
s#lxcfs-agent:.*#lxcfs-agent:v$(IMAGE_TAGS)#
endef

export UPDATE_CHART_YAML
export UPDATE_VALUES_YAML
$(CHARTS_DIRECTORY)/lxcfs-on-kubernetes-$(RELEASE_VERSION).tgz: helm $(CHART_TEMPLATE_PATH)/clusterrole.yaml \
	$(CHART_TEMPLATE_PATH)/certificate.yaml $(CHART_TEMPLATE_PATH)/issuer.yaml \
	$(CHART_TEMPLATE_PATH)/mutatingwebhookconfiguration.yaml \
	$(CHART_TEMPLATE_PATH)/deployment.yaml \
	$(CHART_TEMPLATE_PATH)/daemonset.yaml \
	$(CHART_TEMPLATE_PATH)/clusterrolebinding.yaml

	sed -i.bak -E "$$UPDATE_CHART_YAML" ${CHART_PROJECT_PATH}/Chart.yaml
	sed -i.bak -E "$$UPDATE_VALUES_YAML" ${CHART_PROJECT_PATH}/values.yaml
	$(MAKE) $(CHART_PROJECT_PATH)/README.md
	$(HELM) package $(CHART_PROJECT_PATH) \
		--version $(RELEASE_VERSION) \
		--app-version $(RELEASE_VERSION) \
		--destination $(CHARTS_DIRECTORY)
	-mv ${CHART_PROJECT_PATH}/values.yaml.bak ${CHART_PROJECT_PATH}/values.yaml
	-mv ${CHART_PROJECT_PATH}/Chart.yaml.bak ${CHART_PROJECT_PATH}/Chart.yaml
.PHONY: helm-generate
helm-generate: $(CHARTS_DIRECTORY)/lxcfs-on-kubernetes-$(RELEASE_VERSION).tgz

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)