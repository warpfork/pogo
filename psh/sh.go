package psh

import (
	"fmt"
	"os"
	"os/exec"
)

func Sh(cmd string) sh {
	var cmdt commandTemplate
	cmdt.cmd = cmd
	return func(args ...interface{}) sh {
		return closure(cmdt, args...)
	}
}

type sh func(args ...interface{}) sh

func closure(cmdt commandTemplate, args ...interface{}) sh {
	fmt.Fprintf(os.Stderr, "closure! :: %d args: %#v\n\n", len(args), args)

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "running :: %#v\n\n", cmdt)

		bareCmd := exec.Command(cmdt.cmd, cmdt.args...)
		// set up direct stdin by hack for now
		bareCmd.Stdin = os.Stdin
		bareCmd.Stdout = os.Stdout
		bareCmd.Stderr = os.Stderr
		cmd := NewRunningCommand(bareCmd)
		cmd.Start()
		cmd.Wait()
		return nil
	} else {
		fmt.Fprintf(os.Stderr, "modifying :: %#v\n\n", cmdt)
		for _, rarg := range args {
			switch arg := rarg.(type) {
			case string:
				cmdt.bakeArgs(arg)
			default:
				// ignore, for now
			}
		}
		return func(x ...interface{}) sh {
			return closure(cmdt, x...)
		}
	}
}

// func (sh sh) BakeArgs(opts Opts) sh {
// 	return 
// }

func (cmdt *commandTemplate) bakeArgs(args ...string) {
	cmdt.args = append(cmdt.args, args...)
}
