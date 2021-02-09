SOURCE_DATE := $(shell date --utc +"%Y-%m-%dT%H:%M:%S")
VERSION?=dev
GITHUB_VERSION := $(subst /,-,$(GITHUB_REF)_$(GITHUB_SHA)_$(GITHUB_RUN_ID))

clean:
	@rm -rf ./build

test:
	@go test -v ./...

build: clean
	@$(GOPATH)/bin/goxc \
		-bc="linux,windows" \
		-pv=$(VERSION) \
		-d=build \
		-build-ldflags "-X main.VERSION=$(VERSION) -X main.SOURCE_DATE=$(SOURCE_DATE)"

version:
	@echo $(VERSION)

githubversion:
	@echo dev@$(GITHUB_VERSION)
