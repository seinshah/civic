.PHONY: gomod gomod-upgrade generate lint test vulncheck validate

GO := go
LINTER := golangci-lint

gomod:
	@$(GO) mod tidy
	@$(GO) mod vendor

gomod-upgrade: package := ./...
gomod-upgrade:
	@$(GO) get -u $(package)
	@$(GO) mod tidy
	@$(GO) mod vendor

generate:
	@$(GO) generate ./...

lint:
	@$(LINTER) run --timeout=5m

test:
	@$(GO) test -cover -race -timeout=5m ./...

vulncheck: ## Check the code against recent vulnerabilities
	@govulncheck ./...

validate: generate lint vulncheck test