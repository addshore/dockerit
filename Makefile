VERSION=0.0.1

clean:
	@rm -rf ./build

build: clean
	@$(GOPATH)/bin/goxc \
			-bc="linux,386" \
		-pv=$(VERSION) \
		-d=build \
		-build-ldflags "-X main.VERSION=$(VERSION)"

version:
	@echo $(VERSION)
