.PHONY: tag

tag: version ?=
tag:
ifdef version
	@rm -rf ./versioned_docs/version-$(version) || echo "Skip removing docs for this version..."
	@rm -f ./versioned_sidebars/version-$(version)-sidebars.json || echo "Skip removing sidebar file for this version..."
	@jq 'del(.[] | select(. == "0.1"))' ./versions.json > ./temp-versions.json && mv ./temp-versions.json ./versions.json
	npm run docusaurus docs:version $(version)
else
	@echo "Please provide a version number"
endif
