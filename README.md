# sensu-yieldbot-handler-elasticsearch
ES Handler in golang

[![Build Status][travis-image]][travis-url]

### Building (w/o Vagrant)

These instructions assume you are on an OSX machine and you are already configured for cross-compiling.

The easiest way to configure OSX for cross-compiling is using [Homebrew](http://brew.sh/). Use the following command to have everything configured from the box automatically **highly reccomended** `brew install go --cross-compile-all`. Then all you need to do is follow the below commands as examples. `-o` is the binary to output, I prefix them so I know which is which. You **must** have CGO disabled for go to build cross-platform.

To build for OSX use:
`go build -o binary.osx path/to/file.go`

To build for linux use:
`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o binary.linux /path/to/file.go`

### Building (w/ Vagrant)

`vagrant up`
`vagrant ssh`

`cd /opt/gopath/github.com/yieldbot/diemon`
`make release`

the binary will be in *./bin* and the zip file will be in *./pkg*. Currently this will only build for linux/amd64 but you can call gox with any range of os and arch options. Use `gox --help` for all available choices.

[travis-url]: https://travis-ci.org/yieldbot/diemon
[travis-image]: https://travis-ci.org/yieldbot/diemon.svg?branch=master
