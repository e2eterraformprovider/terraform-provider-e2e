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
	go test -count=1 $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test -count=1 $(TESTARGS) -timeout=60s -parallel=4

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

.PHONY: build test testacc vet fmt fmtcheck lint install _upgrade_goe2e upgrade_goe2e vendor changelog

# =============================================================================
# COLORS FOR OUTPUT
# =============================================================================
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

# =============================================================================
# CHANGELOG AUTOMATION
# =============================================================================
changelog:  ## Update CHANGELOG.md with new version entry from git commits
	@echo "$(BLUE)📝 Generating CHANGELOG entry...$(NC)"
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)❌ VERSION is required. Usage: make changelog VERSION=2.2.8$(NC)"; \
		exit 1; \
	fi
	@if [ ! -f "CHANGELOG.md" ]; then \
		echo "$(RED)❌ CHANGELOG.md not found$(NC)"; \
		exit 1; \
	fi
	@VERSION=$(VERSION); \
	DATE=$$(date +%Y-%m-%d); \
	echo "$(YELLOW)➕ Adding entry for version $$VERSION ($$DATE)$(NC)"; \
	cp CHANGELOG.md CHANGELOG.md.bak; \
	LAST_TAG=$$(git describe --tags --abbrev=0 2>/dev/null || echo ""); \
	if [ -n "$$LAST_TAG" ]; then \
		echo "$(BLUE)📋 Extracting commits since $$LAST_TAG...$(NC)"; \
		COMMITS=$$(git log $$LAST_TAG..HEAD --pretty=format:"- %s" --no-merges 2>/dev/null | grep -v "^- Merge" | head -20); \
		if [ -z "$$COMMITS" ]; then \
			COMMITS="- Version bump"; \
		fi; \
	else \
		echo "$(YELLOW)⚠️  No previous tag found, using recent commits$(NC)"; \
		COMMITS=$$(git log -10 --pretty=format:"- %s" --no-merges 2>/dev/null || echo "- Initial release"); \
	fi; \
	{ \
		head -n 10 CHANGELOG.md | grep -B 10 "^## \[Unreleased\]" || head -n 6; \
		echo ""; \
		echo "## [$$VERSION] - $$DATE"; \
		echo ""; \
		echo "### Changed"; \
		echo "$$COMMITS"; \
		echo ""; \
		echo "---"; \
		echo ""; \
		tail -n +11 CHANGELOG.md | sed '/^## \[Unreleased\]/,/^---/d'; \
	} > CHANGELOG.md.tmp && mv CHANGELOG.md.tmp CHANGELOG.md; \
	echo "$(GREEN)✅ CHANGELOG.md updated with version $$VERSION$(NC)"; \
	echo "$(YELLOW)💡 Review and edit the entry, then commit:$(NC)"; \
	echo "   git diff CHANGELOG.md"; \
	echo "   git add CHANGELOG.md"; \
	echo "   git commit -m \"chore: update CHANGELOG for v$$VERSION\""; \
	echo ""; \
	echo "$(BLUE)📦 Backup saved as CHANGELOG.md.bak$(NC)"

_upgrade_goe2e:
	@echo "==> upgrading goe2e"
	@go get -u github.com/e2enetworks/goe2e
	@echo "==> upgraded goe2e"
	@echo ""

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
