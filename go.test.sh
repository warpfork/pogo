#!/bin/bash

function gotest {
	package="$1"; shift
	GORACE="log_path=test-$package-race.log" GOPATH=$PWD/.go/ go test -race "polydawn.net/gosh/$package" "$@"
}

gotest log
gotest picnic
gotest prom
gotest psh
