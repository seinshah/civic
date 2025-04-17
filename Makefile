CONTAINER_NAME := civic

GO := go
LINTER := golangci-lint
DOCKER := docker
VULNCHECKER := govulncheck
NILCHECKER := nilaway
EXECUTABLE := ./main.go

TARGET_RUNNER =

.PHONY: change-runner
change-runner: runner ?= host
change-runner: ## change the runner to run commands on the host machine or docker - pass runner=<host|docker> argument
ifeq ($(runner), host)
	@sed -i '' 's/^TARGET_RUNNER =.*/TARGET_RUNNER = /' Makefile
else ifeq ($(runner), docker)
	@sed -i '' 's/^TARGET_RUNNER =.*/TARGET_RUNNER = $$(DOCKER) exec $$(CONTAINER_NAME)/' Makefile
else
	@echo "Please provide a valid runner: host or docker"
endif

.PHONY: help
help: ## Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: build
build: ## build the image with dev target if runner is docker
ifneq ($(TARGET_RUNNER),)
ifeq ($(shell $(DOCKER) ps -q -f name=$(CONTAINER_NAME)),)
	@$(DOCKER) build --target dev -t $(CONTAINER_NAME):dev .
	@$(DOCKER) run -d --rm -v $$(pwd):/app --name $(CONTAINER_NAME) $(CONTAINER_NAME):dev
endif
endif

.PHONY: stop
stop: ## stop and remove the container
ifneq ($(shell $(DOCKER) ps -q -f name=$(CONTAINER_NAME)),)
	@$(DOCKER) stop $(CONTAINER_NAME)
endif

.PHONY: rebuild
rebuild: stop build ## rebuild the image and restart the container

.PHONY: cmd-generate-cv
cmd-generate-cv: version := 0.1
cmd-generate-cv: output :=
cmd-generate-cv: schema :=
cmd-generate-cv: ## generate cv with the provided schema to the provided output - args: output, schema
	@$(TARGET_RUNNER) $(GO) run -ldflags "-X main.Version=$(version)" $(EXECUTABLE) generate -s $(schema) -o $(output)

.PHONY: cmd-json-schema
cmd-json-schema: ## generate jsonschema of civic config file in the current directory
	@$(TARGET_RUNNER) $(GO) run -ldflags "-X main.Version=$(version)" $(EXECUTABLE) schema json

.PHONY: cmd-init-schema
cmd-init-schema: output :=
cmd-init-schema: ## generate a civic schema template for you to build upon
	@$(TARGET_RUNNER) $(GO) run -ldflags "-X main.Version=$(version)" $(EXECUTABLE) schema init -o $(output)

.PHONY: cmd-generate-template-example
cmd-generate-template-example: version := 0
cmd-generate-template-example: template :=
cmd-generate-template-example: ## generate example html and pdf outputs for a given template in the given version
	@$(TARGET_RUNNER) sed -E "/^template:/,/^[^[:space:]]/ s|^([[:space:]]*path:[[:space:]]*).*|\1./templates/$(template)/v$(version)/template.html|" ./examples/example.schema.yaml > /tmp/$(template)_sample_schema.yaml
	@$(TARGET_RUNNER) $(GO) run -ldflags "-X main.Version=$(version)" $(EXECUTABLE) generate -s /tmp/$(template)_sample_schema.yaml -o ./templates/$(template)/v$(version)/example.html
	@$(TARGET_RUNNER) $(GO) run -ldflags "-X main.Version=$(version)" $(EXECUTABLE) generate -s /tmp/$(template)_sample_schema.yaml -o ./templates/$(template)/v$(version)/example.pdf
	@$(TARGET_RUNNER) rm -f /tmp/$(template)_sample_schema.yaml

.PHONY: gomod
gomod: build ## clean up imported go packages
	@$(TARGET_RUNNER) $(GO) mod tidy
	@$(TARGET_RUNNER) $(GO) mod vendor

.PHONY: gomod-upgrade
gomod-upgrade: package := ./...
gomod-upgrade: build ## upgrade go modules - use package argument to upgrade specific package
	@$(TARGET_RUNNER) $(GO) get -u $(package)
	@$(TARGET_RUNNER) $(GO) mod tidy
	@$(TARGET_RUNNER) $(GO) mod vendor

.PHONY: validate
validate: generate lint vulncheck test ## validate the code against all the validator in one place

.PHONY: generate
generate: build ## automatically re-generate code
	@$(TARGET_RUNNER) $(GO) generate ./...

.PHONY: lint
lint: build ## check formatting
	@$(TARGET_RUNNER) $(LINTER) run --timeout=5m

.PHONY: test
test: build ## run tests
	@$(TARGET_RUNNER) $(GO) test -cover -race -timeout=5m ./...

.PHONY: vulncheck
vulncheck: build ## Check the code against recent vulnerabilities
	@$(TARGET_RUNNER) $(VULNCHECKER) ./...

.PHONY: nilcheck
nilcheck: build ## Check nil pointer dereference issues
	@$(TARGET_RUNNER) $(NILCHECKER) ./...

.PHONY: release
release: version := '0.1.0-dev'
release: platforms := 'linux/arm64'
release: ## Create the production grade image of the app with the provided version tag and platforms
	@echo "Building v$(version) on $(platforms)...\n"
	@$(DOCKER) buildx build . \
			--platform=$(platforms) \
			--build-arg="APP_VERSION=$(version)" \
			--tag=$(CONTAINER_NAME):latest \
			--tag=$(CONTAINER_NAME):v$(version)
