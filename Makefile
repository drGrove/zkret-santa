PACKAGE_NAME = zkret-santa
IMAGE_NAME_STANDARD = secret-santa
IMAGE_NAME_ZKSNARK = $(PACKAGE_NAME)

# Container build configuration
PLATFORM ?= linux/amd64,linux/arm64,linux/arm/v7
PROGRESS ?= auto
SOURCE_DATE_EPOCH := $(shell git log -1 --format=%ct 2>/dev/null || echo 0)
COMMIT_ISO := $(shell git log -1 --format=%cI 2>/dev/null || echo "unknown")
REGISTRY ?= ghcr.io/drgrove
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
OUT ?= out
IMAGE_DIR_STANDARD := $(OUT)/$(IMAGE_NAME_STANDARD)
IMAGE_DIR_ZKSNARK := $(OUT)/$(IMAGE_NAME_ZKSNARK)
BINARY_NAME ?= $(PACKAGE_NAME)
GO_SRCS := $(wildcard **/*.go) $(wildcard go.*) Makefile

export SOURCE_DATE_EPOCH
export TZ=UTC
export LANG=C.UTF-8
export LC_ALL=C
export BUILDKIT_MULTI_PLATFORM=1
export DOCKER_BUILDKIT=1

.PHONY: build
build: $(GO_SRCS)
	go build \
		$(EXTRA_ARGS) \
		-trimpath \
		-v \
		-mod=readonly \
		-o ${BINARY_NAME} \
		./cmd/$(PACKAGE_NAME)

.PHONY: install
install: $(GO_SRCS)
	go install \
		$(EXTRA_ARGS) \
		./cmd/$(PACKAGE_NAME)

.PHONY: test
test: $(GO_SRCS)
	go test \
		-v \
		$(EXTRA_ARGS) \
		./...

$(IMAGE_NAME_STANDARD): $(GO_SRCS)
	$(MAKE) build BINARY_NAME=$@

.PHONY: test-$(IMAGE_NAME_STANDARD)
test-$(IMAGE_NAME_STANDARD): test

.PHONY: install-$(IMAGE_NAME_STANDARD)
install-$(IMAGE_NAME_STANDARD):
	go install ./cmd/$(PACKAGE_NAME)

.PHONY: run
run: build
	./$(PACKAGE_NAME)

# zkSNARK-enabled build targets
$(PACKAGE_NAME): $(GO_SRCS)
	$(MAKE) build EXTRA_ARGS="-tags zksnark"

.PHONY: test
test-$(IMAGE_NAME_ZKSNARK):
	go test -tags zksnark -v ./...

.PHONY: install
install-$(IMAGE_NAME_ZKSNARK):
	$(MAKE) install EXTRA_ARGS="-tags zksnark"

.PHONY: clean
clean:
	rm -f $(IMAGE_NAME_STANDARD) $(PACKAGE_NAME)
	rm -rf $(OUT)

$(OUT):
	@mkdir -p $@

$(OUT)/%/index.json: Containerfile | $(OUT)
	$(eval STEM := $*)
	$(eval BUILD_TAGS := $(if $(filter zkret-santa,$(STEM)),zksnark,))
	# TODO: remove this when llvm and arm64 merges in stagex
	$(eval PLATFORM := $(if $(filter zkret-santa,$(STEM)),amd64,$(PLATFORM)))
	@mkdir -p $(OUT)/$(STEM)
	docker buildx build \
		--ulimit nofile=2048:16384 \
		--tag $(REGISTRY)/$(STEM):$(VERSION) \
		--build-arg BUILD_TAGS=$(BUILD_TAGS) \
		--build-arg BINARY_NAME=$(STEM) \
		--target package-$(STEM) \
		--output type=oci,rewrite-timestamp=true,force-compression=true,tar=true,dest=- \
		--platform=$(PLATFORM) \
		--progress=$(PROGRESS) \
		--sbom=true \
		--provenance=true \
		-f Containerfile \
		. | tar -C $(OUT)/$(STEM) -xm

image-%:
	$(MAKE) $(OUT)/$*/index.json

$(OUT)/%:
	@mkdir -p $(OUT)

# Build both images
.PHONY: images
images: image-$(IMAGE_NAME_STANDARD) image-$(IMAGE_NAME_ZKSNARK)

# Digest extraction - create digest files for verification
.PHONY: image-digests-%
image-digests-%:
	$(MAKE) $(OUT)/$*.txt

$(OUT)/%.txt: $(OUT)/%/index.json
	@INDEX_DIGEST=$$(jq -r '.manifests[0].digest' $(OUT)/$*/index.json) && \
	MANIFEST_FILE=$$(echo "$$INDEX_DIGEST" | sed 's/sha256://' | xargs -I {} find $(OUT)/$*/blobs/sha256 -name "{}" -type f) && \
	if [ -n "$$MANIFEST_FILE" ]; then \
		jq -r '.manifests[] | select(.annotations."vnd.docker.reference.type" != "attestation-manifest") | "\(.digest | sub("sha256:"; "")) \(.platform.os)/\(.platform.architecture)" + if .platform.variant then "/" + .platform.variant else "" end' "$$MANIFEST_FILE" | sort > $@; \
	else \
		echo "Error: Could not find manifest file for $$INDEX_DIGEST"; \
		exit 1; \
	fi

image-digests: image-digests-$(IMAGE_NAME_STANDARD) image-digests-$(IMAGE_NAME_ZKSNARK)

show-digests-%: image-digests-%
	@cat $(OUT)/$*.txt

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'

