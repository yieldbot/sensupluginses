SHELL = /bin/sh
.PHONY: all build bump_version clean coverage dist format help install lint maintainer-clean test test_all updatedeps version vet

# We only care about golang and texinfo files at the moment so clear and explictly denote that
.SUFFIXES:
.SUFFIXES: .go .texinfo

# Set the location for installing GNU info files
# You can overwrite this by setting your build command to
# `make infodir=path install`
ifndef infodir
infodir = /usr/local/share/info
endif

# Set the src directory. You can overwrite this by setting your build command
# to `make srcdir=path build`
ifndef srcdir
srcdir = src
endif

# Set the default os/arch to build for. Specify additional values in a space
# seperated array. To overwrite this use
# `make osarch="linux/amd64 linux/386" build`
ifndef osarch
osarch = linux/amd64
endif

# Set the base package location.
# `make pkgbase="yieldbot" build
ifndef pkgbase
pkgbase = github.com
endif

# Set the default package. Specify additional values in a space
# seperated array. To overwrite this use
# `make pkg="diemon bobogono" build`
ifndef pkg
pkg := $(shell pwd | awk -F/ '{ print $$NF }')
endif

# Set the name of the output file. If using only a single os/arch the name
# will be as given. If multiple os/arch combinations are used then the given
# name will be suffixed with _OS_ARCH.
ifndef out
	ifeq ("$(osarch)","linux/amd64")
		output = ./bin/$(pkg)
	else
			output = ./bin/$(pkg)_{{.OS}}_{{.Arch}}
	endif
else
	ifeq ("$(osarch)","linux/amd64")
		output = ./bin/$(out)
	else
			output = ./bin/$(out)_{{.OS}}_{{.Arch}}
	endif
endif

# Set the path that the tarball will be dropped into. DrTeeth will look in
# ./target by default but golang will put it into ./pkg if left to itself.
ifndef target_path
target_path = target
endif

# Set where the local binary should be installed to for testing purposes.
ifndef destdir
destdir = /usr/local/bin
endif

default: all

# build and then create a tarball in the target directory
# basically everything needed to put it into artifactory
all: format lint updatedeps build dist

# Build a binary from the given package and drop it into the local bin
build:
	@for i in $$(echo $(pkg)); do \
  	gox -osarch="$(osarch)" -output=$(output) $(pkgbase)/$$i/$(srcdir); \
  done; \

# bump the version of the project
bump_version:
	@ver=$$(awk '{ print $$NF }' version | awk -F. '{ print $$NF }'); \
	ver=$$(($$ver+1)); \
	echo "version 0.0.$$ver" > ./version


# delete all files used for building
clean:
	@echo "this needs to be implemented"

# run the golang coverage tool
coverage:
	@echo "this needs to be implemented"

# pack everything up neatly
dist: build
	tar -czvpf $(target_path)/$(pkg).tgz $$GOPATH/src/$(pkgbase)/$(pkg)/bin/*

# run the golang formatting tool on all files in the current src directory
format:
	@gofmt -w $(srcdir)/*.go

# install the binary and any info docs locally for testing
install:
	@if [ -e $$GOPATH/src/$(pkgbase)/$(pkg)/bin/* ]; then \
	  mkdir -p $(destdir); \
	  cp $$GOPATH/src/$(pkgbase)/$(pkg)/bin/* $(dest dir); \
	else \
		echo "Nothing to install, no binaries were found in ./bin/"; \
	fi; \

	@if [ -e ./docs/*.info ]; then \
	  mkdir -p $(infodir); \
	  cp ./docs/*.info $(infodir); \
	fi; \

info:
	@if [ -e ./docs/*.texinfo ]; then \
	  makeinfo ./docs/*.texinfo --output ./docs/; \
	else \
		echo "Nothing to build, no info files were found in ./docs/"; \
	fi; \

# run the golang linting tool
lint:
	@golint $$GOPATH/src/$(pkgbase)/$(pkg)/$(srcdir)/*.go

maintainer-clean:
	@echo "this needs to be implemented"

# run unit tests and anything else testing wise needed
test:
	@echo "this needs to be implemented"

# run unit tests, vet, formatting, and liniting combined
test_all:
	@echo "this needs to be implemented"

# update all deps to the latest versions available
updatedeps:
	@go list ./... \
		| xargs go list -f '{{join .Deps "\n"}}' \
		| sort -u \
		| xargs go get -f -u -v

# print out the current version of the project
version:
	@if [ -e version ]; then \
	  awk '{ print $$NF }' version; \
	else \
		@echo "No version file found in the project root"; \
	fi; \

# run go vet
vet:
	@echo "this needs to be implemented"
