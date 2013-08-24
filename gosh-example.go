package main

import (
	"polydawn.net/gosh/psh"
)

func main() {
	psh.Command("echo").
	Bake("word1").
	Bake("word2").
	Go()

	psh.Command("bash").
	Bake("-c").
	Bake("exit 44").
	Go()

	troll("bash")("wat")

	troll(
		"bash",
		)(
		"wat",
	)

	troll("bash",)("wat")

	troll("bash",
	)("wat")

	troll("bash",
	)("wat",
	)()

	troll("bash")(
	"wat")(
	)

	troll("echo",
	)("-ne",
	)(
		"wat",
		"bat",
	)()

	echo := troll("echo")("-ne")
	echo("wat", "bat")()
	echo("wat", "bat")(psh.Opts{In:"zat"})()

	troll()().megaboop()()()().megaboop()
}

type trollr func(hax ...interface{}) trollr

func troll(hax ...interface{}) trollr {
	return troll
}

// if you build the deeper layer around trollr, you can build the upper layer as structs with a very small set of function names so you can have your choice of syntax styles.
// which seems like a good thing to do seeing as how go is kind of crazy about letting open-paren start a new line, i mean.

// func troll(hax ...interface{}) func(...interface{}) func(...interface{}) {
// 	return nil
// }

//func (hax *func(int)) Boop() {}

func (hax trollr) megaboop() trollr { return hax }
