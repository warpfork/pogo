package main

import (
	. "polydawn.net/pogo/gosh"
	"os"
)

func main() {
	echo := Sh("echo")("-n", "-e").
		BakeOpts(DefaultIO)

	echo("wat\n", "\t\033[0;31mred and indented\033[0m\n")()

	cat := Sh("cat")
	catIn := cat.BakeArgs("-")
	catIn(Opts{In:"piping in things is easy!\n", Out: os.Stdout})()

	shell := Sh("bash")("-c")

	shell("echo $TERM > testlag")()

	shell(ClearEnv{})("echo $TERM > testlag2")()

	shell(Env{"VAR": "59"})(Opts{OkExit:[]int{59}})("exit $VAR")()
}
