#!/bin/bash

export GOPATH="$PWD"/.go/

function gotest {
	package="$1"; shift
	GORACE="log_path=test-$package-race.log" go test -race "polydawn.net/gosh/$package" "$@"
}

gotest log
gotest picnic
gotest prom
gotest psh
