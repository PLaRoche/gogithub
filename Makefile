VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
GO_BUILD=CGO_ENABLED=0 go build -i --ldflags="-w"

LINTERS=\
	gofmt \
	golint \
	vet \
	misspell \
	ineffassign \
	deadcode

all: ci
ci: $(LINTERS) build

.PHONY: all ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
	github.com/golang/lint/golint \
	honnef.co/go/simple/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter \
	github.com/jteeuwen/go-bindata/... \
	golang.org/x/oauth2 \
	github.com/google/go-github/github

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	glide -v || curl http://glide.sh/get | sh

vendor: glide.lock
	glide install

# this target is used for pre-building all packages, allowing us to cache them
# during docker-compose builds.
pkginstall:
	$(GO_BUILD) $$(glide nv)

.PHONY: bootstrap $(BOOTSTRAP)
.PHONY: pkginstall

#################################################
# Test and linting
#################################################

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	$$(glide nv | grep -v generated) --enable

$(LINTERS): vendor
	$(METALINT) $@

.PHONY: $(LINTERS)

#################################################
# Building
#################################################

build: bin/gogithub

bin/gogithub: vendor gogithub.go
	$(GO_BUILD) -o bin/gogithub

.PHONY: build

#################################################
# Cleaning
#################################################

clean:
	rm bin/gogithub
