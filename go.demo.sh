#!/bin/bash

export GOPATH="$PWD"/.gopath/

go build -race gosh-example.go && ./gosh-example
