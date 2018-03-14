.PHONY: setup
setup:
	$(MAKE) install-requirements
	glide cc && glide install --force -v

.PHONY: install-requirements
install-requirements:
	@type glide >/dev/null 2>&1 || curl https://glide.sh/get | sh
	go get -t -u github.com/goreleaser/goreleaser

.PHONY: release
release:
	goreleaser --skip-publish --rm-dist
