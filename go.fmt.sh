#!/bin/bash

export GOPATH="$PWD"/.gopath/

function gothing {
	package="$1"; shift
	go fmt "polydawn.net/pogo/$package" "$@"
}

gothing log
gothing picnic
gothing iox
gothing prom
gothing psh
