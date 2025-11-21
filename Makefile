TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=e2eterraformprovider
NAME=e2e
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

lint:
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run ./...

install-terrafmt:
	@echo "==> Installing terrafmt..."
	@go install github.com/katbyte/terrafmt@latest

terrafmt-check: install-terrafmt
	@echo "==> Checking terraform code with terrafmt..."
	@if [ -d "docs" ]; then \
		terrafmt diff --check --fmtcompat docs/ || (echo "Terraform code blocks in docs/ are not formatted. Run 'make terrafmt' to fix."; exit 1); \
	fi
	@if [ -d "examples" ]; then \
		terrafmt diff --check --fmtcompat examples/ || (echo "Terraform code blocks in examples/ are not formatted. Run 'make terrafmt' to fix."; exit 1); \
	fi

terrafmt: install-terrafmt
	@echo "==> Formatting terraform code with terrafmt..."
	@if [ -d "docs" ]; then \
		terrafmt fmt docs/; \
	fi
	@if [ -d "examples" ]; then \
		terrafmt fmt examples/; \
	fi

fmt:
	gofmt -w -s .

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=all $(SWEEPARGS) -timeout 60m

.PHONY: build install test testacc vet lint install-terrafmt terrafmt-check terrafmt fmt sweep
