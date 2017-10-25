### Makefile for killy
GITTAG=`git rev-parse --short HEAD`
BUILD_TIME=`date -u +%Y.%m.%d-%H:%M:%S%Z`
VERSION=0.0.1
GOPATH ?= $(shell go env GOPATH)
GOFLAGS=-ldflags "-X "github.com//killy/internal/version".GitCommit=${GITTAG} -X "github.com/hawkingrei/killy/internal/version".BuildTime=${BUILD_TIME} -X "github.com/hawkingrei/killy/internal/version".Version=${VERSION}"

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif
CURDIR := $(shell pwd)
path_to_add := $(addsuffix /bin,$(subst :,/bin:,$(CURDIR)/_vendor:$(GOPATH)))
export PATH := $(path_to_add):$(PATH)

FILES     := $$(find . -name '*.go' | grep -vE 'vendor')
BLDDIR    := build

APPS := killyd

all: $(APPS)

linux :
	$(wildcard apps/killyd/*.go       killyd/*.go       killy/*.go internal/*/*.go)	
	@rm -fr $(BLDDIR)
	@mkdir -p $(dir $@)
	GOOS=linux  GOARCH=amd64 go build ${GOFLAGS} -o $@ ./apps/$*

$(BLDDIR)/killyd:        $(wildcard apps/killyd/*.go       killyd/*.go       killy/*.go internal/*/*.go)



$(BLDDIR)/%:
	@rm -fr $(BLDDIR)
	@mkdir -p $(dir $@)
	go build ${GOFLAGS} -o $@ ./apps/$*

$(APPS): %: $(BLDDIR)/%
	
clean:
	@rm -fr $(BLDDIR)

check:
#go get github.com/golang/lint/golint

	@echo "gofmt (simplify)"
	@gofmt -s -l -w $(FILES) 2>&1 | grep -v "parser/parser.go" | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "vet"
	@ go tool vet $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "vet --shadow"
	@ go tool vet --shadow $(FILES) 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'
	@echo "golint"
	@ golint ./... 2>&1 | grep -vE 'vendor|context\.Context|LastInsertId|NewLexer|\.pb\.go' | awk '{print} END{if(NR>0) {exit 1}}'
