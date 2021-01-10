BUILD_DATE := $(shell date --utc +"%Y%m%dT%H%M%S")
BUILD_VERSION?=dev@$(BUILD_DATE)
GITHUB_VERSION := $(subst /,-,$(GITHUB_REF)_$(GITHUB_SHA)_$(GITHUB_RUN_ID)@$(BUILD_DATE))

clean:
	@rm -rf ./build

build: clean
	@$(GOPATH)/bin/goxc \
			-bc="linux,386" \
		-pv=$(BUILD_VERSION) \
		-d=build \
		-build-ldflags "-X main.VERSION=$(BUILD_VERSION)"

version:
	@echo $(BUILD_VERSION)

githubversion:
	@echo dev@$(GITHUB_VERSION)
