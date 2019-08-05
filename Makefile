default: all

GO_PACKAGES = $$(go list ./... | grep -v vendor)
GO_FILES = $$(find . -name "*.go" | grep -v vendor | uniq)

cli-linux:
	GOOS=linux GOARCH=amd64 go build -o ksm.linux ./cmd/cli/main.go

cli-linux32:
	GOOS=linux GOARCH=386 go build -o ksm.linux32 ./cmd/cli/main.go

cli-mac:
	GOOS=darwin GOARCH=amd64 go build -o ksm.darwin ./cmd/cli/main.go

kibosh-linux:
	GOOS=linux GOARCH=amd64 go build -o kibosh.linux ./cmd/kibosh/main.go

kibosh-mac:
	GOOS=darwin GOARCH=amd64 go build -o kibosh.mac ./cmd/kibosh/main.go

ksm-linux:
	GOOS=linux GOARCH=amd64 go build -o ksmd.linux ./cmd/ksmd/main.go

ksm-mac:
	GOOS=darwin GOARCH=amd64 go build -o ksmd.mac ./cmd/ksmd/main.go

uaa-crud:
	GOOS=linux GOARCH=amd64 go build -o uaa-crud.linux ./main.go

build: uaa-crud

unit-test:
	@go test ${GO_PACKAGES}

fmt:
	goimports -l -w $(GO_FILES)

vet:
	@go vet ${GO_PACKAGES}

test: generate unit-test vet

generate:
	go generate ./...

cleandep:
	rm -rf vendor
	rm -f Gopkg.lock

HAS_DEP := $(shell command -v dep;)
HAS_GOIMPORTS := $(shell command -v dep;)

boostrap: bootstrap

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOIMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
	go mod vendor


all: fmt test build

