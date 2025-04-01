GO ?= go
PKG =./cmd/punch/

GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6
GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@latest

punch:
	$(GO) build ./cmd/punch/

run-ast-explorer:
	bash ./scripts/build_wasm.sh
	$(GO) run ./tools/ast_explorer/

clean:
	rm punch

.PHONY: deps-tools
deps-tools: ## install tool dependencies
	$(GO) install $(GOLANGCI_LINT_PACKAGE)
	$(GO) install $(GOFUMPT_PACKAGE)

.PHONY: format
format: ## checks formatting on backend code
	gofumpt -l -w $(PKG)

.PHONY: lint
lint: ## lints go code
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run $(PKG)

