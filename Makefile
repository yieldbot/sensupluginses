SHELL = /bin/sh

# This is a general purpose Makefile for building golang projects
#
# version 0.0.9
# Copyright (c) 2015 Yieldbot

.PHONY: all build bump_version clean coverage dist format info install lint maintainer-clean test test_all updatedeps version vet

# We only care about golang and texinfo files at the moment so clear and explicitly denote that
.SUFFIXES:
.SUFFIXES: .go .texinfo

# Set the location for installing GNU info files
# You can overwrite this by setting your build command to
# `make infodir=path install`
ifndef infodir
infodir = /usr/local/share/info
endif

# Set the package to build. Specify additional values in a space
# separated array. To overwrite this use
# `make pkg="diemon bobogono" build`
ifndef pkg
pkg = "."
endif

# Set the src directory. You can overwrite this by setting your build command
# to `make srcdir=path build`
ifndef srcdir
srcdir = cmd
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

# Set the repo to look for the package in. Specify additional values in a space
# separated array. To overwrite this use
# `make repo="diemon bobogono" build`
ifndef repo
repo := $(shell pwd | awk -F/ '{ print $$NF }')
endif

# Set the name of the output file. If using only a single os/arch the name
# will be as given. If multiple os/arch combinations are used then the given
# name will be suffixed with _OS_ARCH.
ifndef out
	ifeq ("$(osarch)","linux/amd64")
		output = ./bin/$(pkg)/$(pkg)
	else
			output = ./bin/$(pkg)/$(pkg)_{{.OS}}_{{.Arch}}
	endif
else
	ifeq ("$(osarch)","linux/amd64")
		output = ./bin/$(pkg)/$(out)
	else
			output = ./bin/$(pkg)/$(out)_{{.OS}}_{{.Arch}}
	endif
endif

# Set the path that the tarball will be dropped into. DrTeeth will look in
# ./target by default but golang will put it into ./pkg if left to itself.
ifndef targetdir
targetdir = pkg
endif

# Set where the local binary should be installed to for testing purposes.
ifndef destdir
destdir = /usr/local/bin
endif

define help
--Targets--

all: Attempt to run gofmt golint and if those pass then it will pull in the latest
     dependencies and build the requested binaries. If the build completes without
     errors and taball is created and dropped into the targetdir.
		 Ex. `make pkg=<package> all`

build: Run gox with a pre-defined set of options. By default a binary will be built
       for linux/amd64, named the same as the srcdir, and any output will be placed
       in ./bin. Ex. `make pkg=<package> build`

clean: Remove any files that are used or produced during the building and packaging
       steps. This will include the binary and tarball themselves as well as the
       directories that these would get dropped into. `make clean`

coverage: This needs to be implemented.

dist :Create a compressed tar archive of all binary produced during the build steps.
      The tarball will be placed into the directory defined by the <targetdir> make
      variable. Ex. `make pkg=<package> dist`

format: Run the gofmt tool. This will produce a list of files that do not conform
        to the standards set via golang. If you want to attempt to fix these errors
        automatically see the <format_correct> task.

format_correct: Attempt to automatically correct any issues presented via the gofmt
                tool.

install:  Install any binaries and info files into the directories specified by the
         <destdir> and <infodir>. Ex. `make pkg=<package> install`

info:  Build any texinfo documents found. Ex. `make pkg=<package> info`

help:  Display this help message. Ex. `make help`

lint:  Run the golang linting tool. Ex. `make pkg=<package> lint`

maintainer_clean: This needs to be implemented.

pre-build: Ensure that the necessary directories present. This does not need to be
           called by the user.

pre-dist: Ensure that the necessary directories present. This does not need to be
          called by the user.

test: This needs to be implemented.

test-all: Run all optional testing targets.

--Variables--

infodir Set the location for installing GNU info files.
        Default: /usr/local/share/info

pkg Set the package to build. Specify additional values in a space seperated
        array. Ex. `make pkg="diemon bobogono" build`
        Default: .

srcdir Set the src directory.
       Default: src

osarch Set the default os/arch to build for. Specify additional values in a space
       seperated array. Ex. `make osarch="linux/amd64 linux/386" build`
       Default: linux/amd64

pkgbase Set the base package location.
        Ex. `make pkgbase="github.com/yieldbot" build
        Default: github.com

repo Set the repo to look for the package in. Specify additional values in a space
     seperated array. Ex. `make repo="diemon bobogono" build`
     Default: The top level directory

out  Set the name of the output file. If using only a single os/arch the name
     will be as given. If multiple os/arch combinations are used then the given
     name will be suffixed with _OS_ARCH.

target Set the path that the tarball will be dropped into. DrTeeth will look in
       ./target by default but golang will put it into ./pkg if left to itself.
       Default: target

destdir Set where the local binary should be installed to for testing purposes.
        Default: /usr/local/bin

endef

export help

default: all

# build and then create a tarball in the target directory
# basically everything needed to put it into artifactory
all: format build dist

# Build a binary from the given package and drop it into the local bin
build: pre-build
	@for i in $$(echo $(pkg)); do \
	  export PATH=$$PATH:$$GOROOT/bin:$$GOBIN; \
  	gox -parallel=1 -osarch="$(osarch)" -output=$(output) ./$(srcdir)/$$i; \
  done; \

# delete all existing binaries and directories used for building
clean:
		rm -rf $(srcdir)/bin $(srcdir)/$(targetdir)

# run the golang coverage tool
coverage:
	@echo "this needs to be implemented"

# pack everything up neatly
dist: build pre-dist
	@if [ -e $(srcdir)/cmd/$(pkg) ]; then \
    cd $(srcdir)/bin/$(pkg); \
	  tar czvf ../../$(targetdir)/output.tar.gz *; \
	else \
	  echo "No binaries were found. No output package will be created"; \
	fi; \

# run the golang formatting tool on all files in the current src directory
format:
	@OUT=`gofmt -l ./$(srcdir)/$(pkg)/*.go`; if [ "$$OUT" ]; then echo $$OUT; exit 1; fi

# fix any detected formatting issues
format_correct:
	@gofmt -w ./$(srcdir)/$(pkg)/*.go

# install the binary and any info docs locally for testing
install:
	@if [ -e ./bin/* ]; then \
	  mkdir -p $(destdir); \
	  cp ./bin/$(pkg)/* $(destdir); \
	else \
		echo "Nothing to install, no binaries were found in ./bin/"; \
	fi; \

	@if [ -e ./docs/*.info ]; then \
	  mkdir -p $(infodir); \
	  cp ./docs/$(pkg)/*.info $(infodir); \
	fi; \

info:
	@if [ -e ./docs/$(pkg)/*.texinfo ]; then \
	  makeinfo ./docs/$(pkg)/*.texinfo --output ./docs/$(pkg)/; \
	else \
		echo "Nothing to build, no info files were found in ./docs/"; \
	fi; \

help:
	@echo "$$help"


# run the golang linting tool
lint:
	@OUT=`golint ./$(srcdir)/$(pkg)/*.go`; if [ "$$OUT" ]; then echo $$OUT; exit 1; fi

maintainer-clean:
	@echo "this needs to be implemented"

# create a directory to store binaries in
# YELLOW need to account for updated packages
# YELLOW need to set the repo name automatically
pre-build:
	@if [ -e $(srcdir)/cmd/$(pkg) ]; then \
		echo "Ensuring output binary directory exists"; \
		mkdir -p $(srcdir)/bin/$(pkg); \
	else \
	  echo "No binaries were found. No bin directory will be created"; \
	fi; \
	if [ -e $$GOPATH/src/github.com/yieldbot/ybsensuplugin/Makefile ]; then \
	  echo "Correct dependency directory structure already exists, doing nothing"; \
	else \
		echo "Creating proper build environment and dependency directory structure"; \
		mkdir -p $$GOPATH/src/github.com/yieldbot/ybsensuplugin; \
		cp -R * $$GOPATH/src/github.com/yieldbot/ybsensuplugin; \
	fi; \



pre-dist:
	@if [ -e $(srcdir)/cmd/$(pkg) ]; then \
		echo "Ensuring output tarball directory exists"; \
	  mkdir -p $(srcdir)/$(targetdir); \
	else \
	  echo "No binaries were found. No output directory will be created"; \
	fi; \

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
	@if [ -e $(pkg)/version ]; then \
		ver=$$(awk '{ print $$NF }' $(pkg)/version) ;\
    echo "{\"version\":\"$$ver\"}"; \
	else \
		@echo "No version file found in the project root"; \
	fi; \

# bump the version of the project
version_bump:
	@ver=$$(awk '{ print $$NF }' $(pkg)/version | awk -F. '{ print $$NF }'); \
	ver=$$(($$ver+1)); \
	echo "version 0.0.$$ver" > $(pkg)/version

# run go vet
vet:
	@go vet ./$(srcdir)/$(pkg)/*.go
