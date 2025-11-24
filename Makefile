TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME?=e2e
ACCTEST_TIMEOUT?=120m
ACCTEST_PARALLELISM?=2
HOSTNAME=registry.terraform.io
NAMESPACE=e2eterraformprovider
BINARY=terraform-provider-${PKG_NAME}


default: build

build: fmtcheck
	go mod tidy
	@mkdir -p bin
	go build -o bin/$(BINARY)

test: fmtcheck
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test -v ./$(PKG_NAME)/... $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -parallel=$(ACCTEST_PARALLELISM)


vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimports:
	@echo "==> Fixing imports code with goimports..."
	@find . -name '*.go' | grep -v vendor | grep -v generator-resource-id | while read f; do goimports -w "$$f"; done

install-golangci-lint:
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2

lint: install-golangci-lint
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run ./...

fmt:
	gofmt -w -s .

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: build test testacc vet fmt fmtcheck lint install _upgrade_goe2e upgrade_goe2e vendor

_upgrade_goe2e:
#	go get -u github.com/e2enetworks/goe2e
	@echo "==> upgraded goe2e"

upgrade_goe2e: _upgrade_goe2e vendor
	@echo "==> upgrade the goe2e version"
	@echo ""

vendor:
	@echo "==> vendor dependencies"
	@echo ""
	go mod vendor
	go mod tidy

install: build
	@PLATFORM="darwin_arm64"; \
	VERSION="0.1.0"; \
	PLUGIN_DIR=$$HOME/.terraform.d/plugins/registry.terraform.io/e2eterraformprovider/e2e/$$VERSION/$$PLATFORM; \
	echo "Version: $$VERSION | Platform: $$PLATFORM | Installing provider to Terraform plugin directory"; \
	mkdir -p $$PLUGIN_DIR && cp bin/$(BINARY) $$PLUGIN_DIR/; \
	echo "==> Provider installed successfully to $$PLUGIN_DIR"
