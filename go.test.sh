#!/bin/bash

export GOPATH="$PWD"/.gopath/

function gotest {
	package="$1"; shift
	GORACE="log_path=test-$package-race.log" go test -race -v "polydawn.net/pogo/$package" "$@"
}

gotest log
gotest picnic
gotest iox
gotest prom
gotest psh
