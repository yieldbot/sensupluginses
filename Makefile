SHELL = /bin/sh

# This is a general purpose Makefile for building golang projects
#
# version 0.0.14
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

# Set the package to build.
ifndef pkg
  pkg := $(shell pwd | awk -F/ '{print $$NF}')
endif

# set the directory to build. Mostly useful only when building librarys
ifndef srcdir
	srcdir = .
endif

# Set the path that the tarball will be dropped into.
ifndef targetdir
  targetdir = pkg
endif

# Set where the local binary should be installed to.
ifndef destdir
destdir = /usr/local/bin
endif

define help
--Targets--

all: Attempt to run gofmt and if it passes build any binaries or libraries. If the build completes 
	   without errors a taball is created and dropped into the targetdir for a binary. This is
	   the default Yieldbot target for the golang build pipeline. Ex. `make all`

build: Run go build with a pre-defined set of options. By default a binary will be built
       for linux/amd64, named the same as the package, and any output will be placed
       in ./bin. Ex. `make build`

       For more fine grained control over bulding check out `go compile` and go `link`.

clean: Remove any files that are used or produced during the building and packaging
       steps. This will include the binaries and tarballs as well as the
       directories that these would get dropped into. Ex. `make clean`

coverage: This needs to be implemented.

dep_tree: This will call updatedeps first to pull in the latest dependencies. At this
	        point it will remove any previous Godeps tree and replace it.

dist: Create a compressed tar archive of all binary produced during the build steps.
      The tarball will be placed into the directory defined by the <targetdir> make
      variable. Ex. `make dist`

format: Run the gofmt tool. This will produce a list of files that do not conform
        to the standards set via golang. If you want to attempt to fix these errors
        automatically see the <format_correct> task.

format_correct: Attempt to automatically correct any issues presented via the gofmt
                tool.

install:  Install any binaries and info files into the directories specified by the
         <destdir> and <infodir>. Ex. `make install`

info:  Build any texinfo documents found. Ex. `make info`

help:  Display this help message. Ex. `make help`

lint:  Run the golang linting tool. Ex. `make lint`

maintainer_clean: This needs to be implemented.

pre-build: Ensure that the necessary directories are present. This does not need to be
           called by the user.

pre-dist: Ensure that the necessary directories are present. This does not need to be
          called by the user.

test: This needs to be implemented.

test-all: Run all optional testing targets.

--Variables--

infodir Set the location for installing GNU info files.
        Default: /usr/local/share/info

pkg Set the package to build. Ex. `make pkg="bobogono" build`
        Default: current working directory

srcdir Set the directory that the sources are in. Mostly only useful when building libraries or when
       the code does not live at the repo root. 
			 Default: .

target Set the path that the tarball will be dropped into. DrTeeth will look in
       ./target by default but golang will put it into ./pkg if left to itself.
       Default: target

destdir Set where the local binary should be installed to for testing purposes.
        Default: /usr/local/bin

endef

export help
default: all

# build and then create a tarball in the target directory
# basically everything needed to create a release.
all: clean build dist

# Build a binary from the given package and drop it into the local bin
build: 
	  @export PATH=$$PATH:$$GOROOT/bin:$$GOBIN; \
	  if [ -e ./cmd ]; then \
      godep go build -o ./bin/$(pkg) --ldflags "-linkmode external -extldflags '-static'" $(srcdir); \
	  else \
	    godep go build --ldflags "-linkmode external -extldflags '-static'" $(srcdir); \
	  fi; \

# delete all existing binaries and directories used for building
clean:
		@rm -rf ./bin ./pkg

# run the golang coverage tool
coverage:
	@echo "this needs to be implemented"

# create a dependency tree, if one already exists it will destroy it first
dep_tree: updatedeps
	rm -rf Godeps
	godep save
	
# pack everything up neatly
dist: build pre-dist
	@if [ -e ./bin ]; then \
    cd ./bin; \
	  tar czvf ../$(targetdir)/output.tar.gz *; \
	else \
	  echo "No binaries were found. No output package will be created"; \
	fi; \

# run the golang formatting tool on all files in the current src directory
format:
	@OUT=`gofmt -l ./..`; if [ "$$OUT" ]; then echo $$OUT; exit 1; fi

# fix any detected formatting issues
format_correct:
	@gofmt -l -w ./..

# install the binary and any info docs locally for testing
install:
	@if [ -e ./bin/* ]; then \
	  mkdir -p $(destdir); \
	  cp ./bin/* $(destdir); \
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

help:
	@echo "$$help"

# run the golang linting tool
lint:
	@OUT=`golint ./...`; if [ "$$OUT" ]; then echo $$OUT; exit 1; fi

maintainer-clean:
	@echo "this needs to be implemented"

null:
	@echo "move along"

# needed for Jenkins builds due to shared Workspaces
pre-build:
	echo "Creating proper build environment and dependency directory structure"; \
	echo "Creating $$GOPATH/src/github.com/yieldbot/$(pkg)"; \
	mkdir -p $$GOPATH/src/github.com/yieldbot/$(pkg); \
	echo "Copying dependencies from $$(pwd) -> $$GOPATH/src/github.com/yieldbot/$(pkg)"; \
	cp -R ./* $$GOPATH/src/github.com/yieldbot/$(pkg); \

pre-dist:
	@if [ -e ./cmd/ ]; then \
		echo "Ensuring output tarball directory exists"; \
		echo "Creating ./$(targetdir)"; \
	  mkdir -p ./$(targetdir); \
	else \
	  echo "No binaries were found. No output directory will be created"; \
	fi; \

# run unit tests and anything else testing wise needed
test:
	@echo "this needs to be implemented"

# run unit tests, vet, formatting, and liniting combined
test_all: vet lint format test

# update all deps to the latest versions available
updatedeps:
	@go list ./... \
		| xargs go list -f '{{join .Deps "\n"}}' \
		| sort -u \
		| xargs go get -f -u -v

# print out the current version of the project
version:
	@if [ -e ./version ]; then \
		ver=$$(awk '{ print $$NF }' ./version) ;\
    echo "{\"version\":\"$$ver\"}"; \
	else \
		echo "No version file found"; \
	fi; \

version_bump:
	@ver_f=$$(awk -F. '{ print $$1 "."$$2 "." }' version); \
	  ver_p=$$(awk -F. '{ print ++$$NF }' version); \
		echo $$ver_f$$ver_p > version

# run go vet
vet:
	@export PATH=$$PATH:$$GOROOT/bin:$$GOBIN; \
	go vet ./...

