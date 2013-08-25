package main

import (
	. "polydawn.net/gosh/psh"
	"os"
)

func main() {
	echo := Sh("echo")("-n", "-e").
		BakeOpts(Opts{In: os.Stdin, Out: os.Stdout})

	echo("wat\n", "\t\033[0;31mred and indented\033[0m\n")()

	cat := Sh("cat")
	catIn := cat.BakeArgs("-")
	catIn()

	shell := Sh("bash")("-c")

	shell("echo $TERM > testlag")()

	shell(ClearEnv{})("echo $TERM > testlag2")()

	shell(Env{"VAR": "59"})("exit $VAR")()
}
