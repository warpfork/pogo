package psh

import (
	"os"
	"os/exec"
)

func Sh(cmd string) sh {
	var cmdt commandTemplate
	cmdt.cmd = cmd
	return enclose(&cmdt)
}

type sh func(args ...interface{}) sh

// private type, used exactly once to create a const nobody else can create so we can use it as a flag to trigger private behavior
type expose_t bool

const expose expose_t = true

type exposer struct{ cmdt *commandTemplate }

func closure(cmdt commandTemplate, args ...interface{}) sh {
	if len(args) == 0 {
		// an empty call is a trigger for actually starting execution.
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
		// produce a function that when called with an exposer, exposes its cmdt.
		return func(x ...interface{}) sh {
			t := x[0].(*exposer)
			t.cmdt = &cmdt
			return nil
		}
	} else {
		// examine each of the arguments, modify our (already forked) cmdt, and
		//  return a new callable sh closure with the newly baked command template.
		for _, rarg := range args {
			switch arg := rarg.(type) {
			case string:
				cmdt.bakeArgs(arg)
			default:
				panic(IncomprehensibleCommandModifier{wat:&rarg})
			}
		}
		return enclose(&cmdt)
	}
}

func (f sh) expose() commandTemplate {
	var t exposer
	f(expose)(&t)
	return *t.cmdt
}

func enclose(cmdt *commandTemplate) sh {
	return func(x ...interface{}) sh {
		return closure(*cmdt, x...)
	}
}

func (f sh) BakeArgs(args ...string) sh {
	cmdt := f.expose()
	cmdt.bakeArgs(args...)
	return enclose(&cmdt)
}

func (cmdt *commandTemplate) bakeArgs(args ...string) {
	cmdt.args = append(cmdt.args, args...)
}
