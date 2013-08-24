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

// private type, used exactly once to create a const nobody else can create so we can use it as a flag to trigger private behavior
type expose_t bool

const expose expose_t = true

type exposer struct{ cmdt *commandTemplate }

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
	} else if args[0] == expose {
		fmt.Fprintf(os.Stderr, "exposing :: %#v\n\n", cmdt)
		// produce a function that when called exposes its cmdt.
		return func(x ...interface{}) sh {
			t := x[0].(*exposer)
			t.cmdt = &cmdt
			return nil
		}
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

func (f sh) expose() commandTemplate {
	var t exposer
	f(expose)(&t)
	return *t.cmdt
}

func (f sh) BakeArgs(args ...string) sh {
	cmdt := f.expose()
	cmdt.bakeArgs(args...)
	return func(x ...interface{}) sh {
		return closure(cmdt, x...)
	}
}

func (cmdt *commandTemplate) bakeArgs(args ...string) {
	cmdt.args = append(cmdt.args, args...)
}
