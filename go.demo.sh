#!/bin/bash

export GOPATH="$PWD"/.gopath/

go build -race gosh-demo.go \
&& echo -e "build complete!\n--------\n" 1>&2 \
&& ./gosh-demo
