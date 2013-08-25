package psh

import (
	"fmt"
	"os/exec"
	"polydawn.net/gosh/iox"
)

func Sh(cmd string) sh {
	var cmdt commandTemplate
	cmdt.cmd = cmd
	cmdt.env = getOsEnv()
	return enclose(&cmdt)
}

type sh func(args ...interface{}) sh

// private type, used exactly once to create a const nobody else can create so we can use it as a flag to trigger private behavior
type expose_t bool

const expose expose_t = true

type exposer struct{ cmdt *commandTemplate }

func closure(cmdt commandTemplate, args ...interface{}) sh {
	if len(args) == 0 {
		// an empty call is a synonym for sh.Run().
		// if you want to just get a RunningCommand reference to track, use sh.Start() instead.
		enclose(&cmdt).Run()
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
			case Env:
				cmdt.bakeEnv(arg)
			case ClearEnv:
				cmdt.clearEnv()
			case Opts:
				cmdt.bakeOpts(arg)
			default:
				panic(IncomprehensibleCommandModifier{wat: &rarg})
			}
		}
		return enclose(&cmdt)
	}
}

func (f sh) expose() *commandTemplate {
	var t exposer
	f(expose)(&t)
	return t.cmdt
}

func enclose(cmdt *commandTemplate) sh {
	return func(x ...interface{}) sh {
		return closure(*cmdt, x...)
	}
}

func (f sh) BakeArgs(args ...string) sh {
	return enclose(f.expose().bakeArgs(args...))
}

func (cmdt *commandTemplate) bakeArgs(args ...string) *commandTemplate {
	cmdt.args = append(cmdt.args, args...)
	return cmdt
}

func (f sh) BakeEnv(args Env) sh {
	return enclose(f.expose().bakeEnv(args))
}

func (cmdt *commandTemplate) bakeEnv(args Env) *commandTemplate {
	for k, v := range args {
		if v == "" {
			delete(cmdt.env, k)
		} else {
			cmdt.env[k] = v
		}
	}
	return cmdt
}

func (f sh) ClearEnv() sh {
	return enclose(f.expose().clearEnv())
}

func (cmdt *commandTemplate) clearEnv() *commandTemplate {
	cmdt.env = make(map[string]string)
	return cmdt
}

func (f sh) BakeOpts(args ...Opts) sh {
	return enclose(f.expose().bakeOpts(args...))
}

func (cmdt *commandTemplate) bakeOpts(args ...Opts) *commandTemplate {
	for _, arg := range args {
		if arg.Cwd != nil {
			cmdt.Cwd = arg.Cwd
		}
		if arg.In != nil {
			cmdt.In = arg.In
		}
		if arg.Out != nil {
			cmdt.Out = arg.Out
		}
		if arg.Err != nil {
			cmdt.Err = arg.Err
		}
	}
	return cmdt
}

/**
 * Starts execution of the command.  Returns a reference to a RunningCommand,
 * which can be used to track execution of the command, configure exit listeners,
 * etc.
 */
func (f sh) Start() *RunningCommand {
	cmdt := f.expose()
	rcmd := exec.Command(cmdt.cmd, cmdt.args...)

	// set up env
	if cmdt.env != nil {
		rcmd.Env = make([]string, len(cmdt.env))
		i := 0
		for k, v := range cmdt.env {
			rcmd.Env[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}
	}

	// set up opts (cwd/stdin/stdout/stderr)
	if cmdt.Cwd != nil {
		rcmd.Dir = *cmdt.Cwd
	}
	if cmdt.In != nil {
		switch in := cmdt.In.(type) {
		case sh:
			//TODO something marvelous
			panic(fmt.Errorf("not yet implemented"))
		default:
			inreader := iox.ReaderFromInterface(in)
			//TODO if this is nil, raise a (typed) panic (which we don't have a type for yet)
			rcmd.Stdin = inreader
		}
	}
	if cmdt.Out != nil {
		out := iox.WriterFromInterface(cmdt.Out)
		//TODO if this is nil, raise a (typed) panic (which we don't have a type for yet)
		rcmd.Stdout = out
	}
	if cmdt.Err != nil {
		if cmdt.Err == cmdt.Out {
			rcmd.Stderr = rcmd.Stdout
		} else {
			out := iox.WriterFromInterface(cmdt.Err)
			//TODO if this is nil, raise a (typed) panic (which we don't have a type for yet)
			rcmd.Stderr = out
		}
	}

	// go time
	cmd := NewRunningCommand(rcmd)
	cmd.Start()
	return cmd
}

func (f sh) Run() {
	cmdt := f.expose()
	cmd := f.Start()
	cmd.Wait()
	exitCode := cmd.GetExitCode()
	//TODO support configurable expected exit codes
	if exitCode != 0 {
		panic(FailureExitCode{cmdname: cmdt.cmd, code: exitCode})
	}
}
