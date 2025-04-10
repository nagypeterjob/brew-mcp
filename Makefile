GIT_COMMIT:=$(shell git rev-parse --short HEAD)
GIT_TAG:=$(shell git describe --tags --abbrev=0 --exact-match 2> /dev/null)
TAG:=$(if $(GIT_TAG),$(GIT_TAG),build-$(GIT_COMMIT))

.PHONY: build
build:
	go build -o brew-mcp-server -ldflags="-B gobuildid -X 'main.BuildVersion=$(TAG)'" ./cmd
