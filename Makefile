PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

GOOS := $(shell go env GOOS)
GOOSALT ?= 'linux'
ifeq ($(GOOS),'darwin')
  GOOSALT = 'mac'
endif

TAG ?= $(USER)
BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

checks:
	@echo "Checking dependencies"
	@(env bash $(PWD)/buildscripts/checkdeps.sh)

getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golint 1>/dev/null || (echo "Installing golint" && go get -u golang.org/x/lint/golint)
	@which staticcheck 1>/dev/null || (echo "Installing staticcheck" && wget --quiet -O ${GOPATH}/bin/staticcheck https://github.com/dominikh/go-tools/releases/download/2019.1/staticcheck_${GOOS}_amd64 && chmod +x ${GOPATH}/bin/staticcheck)
	@which misspell 1>/dev/null || (echo "Installing misspell" && wget --quiet https://github.com/client9/misspell/releases/download/v0.3.4/misspell_0.3.4_${GOOSALT}_64bit.tar.gz && tar xf misspell_0.3.4_${GOOSALT}_64bit.tar.gz && mv misspell ${GOPATH}/bin/misspell && chmod +x ${GOPATH}/bin/misspell && rm -f misspell_0.3.4_${GOOSALT}_64bit.tar.gz)

crosscompile:
	@(env bash $(PWD)/buildscripts/cross-compile.sh)

verifiers: getdeps vet fmt lint staticcheck spelling

vet:
	@echo "Running $@"
	@GO111MODULE=on go vet github.com/minio/minio/...

fmt:
	@echo "Running $@"
	@GO111MODULE=on gofmt -w -d cmd/
	@GO111MODULE=on gofmt -w -d pkg/
	@GO111MODULE=on gofmt -w -d internal/
	@GO111MODULE=on gofmt -w -d test/


vendor:
	@echo "Running $@"
	@GO111MODULE=on go mod vendor

init:
	@echo "Running $@"
	@GO111MODULE=on go mod init
	@GO111MODULE=on go mod vendor

genproto:
	@echo "Generating protos"
	@(env bash $(PWD)/scripts/genproto.sh)

run:
	@echo "Running $@"
	@GO111MODULE=on go run main.go



# Builds minio, runs the verifiers then runs the tests.
check: test

test:
	@echo "Running integration tests"
	@GO111MODULE=on go test -v -count=1 ./test

verify: build
	@echo "Verifying build"
	@(env bash $(PWD)/buildscripts/verify-build.sh)

coverage: build
	@echo "Running all coverage for minio"
	@(env bash $(PWD)/buildscripts/go-coverage.sh)

# Builds minio locally.
build: checks
	@echo "Building minio binary to './minio'"
	@GO111MODULE=on GOFLAGS="" CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS) -o $(PWD)/minio 1>/dev/null
	@GO111MODULE=on GOFLAGS="" CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS) -o $(PWD)/dockerscripts/healthcheck $(PWD)/dockerscripts/healthcheck.go 1>/dev/null
	@GO111MODULE=on GOFLAGS="" CGO_ENABLED=0 go build -tags kqueue --ldflags $(BUILD_LDFLAGS) -o $(PWD)/dockerscripts/check-user $(PWD)/dockerscripts/check-user.go 1>/dev/null

docker: build
	@docker build -t $(TAG) . -f Dockerfile.dev

# Builds minio and installs it to $GOPATH/bin.
install: build
	@echo "Installing minio binary to '$(GOPATH)/bin/minio'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/minio $(GOPATH)/bin/minio
	@echo "Installation successful. To learn more, try \"minio --help\"."



clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf minio
	@rm -rvf build
	@rm -rvf release

.PHONY: test clean vendor