package psh

import (
	"fmt"
	"os/exec"
	"polydawn.net/gosh/iox"
)

func Sh(cmd string) Shfn {
	var cmdt commandTemplate
	cmdt.cmd = cmd
	cmdt.env = getOsEnv()
	cmdt.OkExit = []int{0}
	return enclose(&cmdt)
}

type Shfn func(args ...interface{}) Shfn

// private type, used exactly once to create a const nobody else can create so we can use it as a flag to trigger private behavior
type expose_t bool

const expose expose_t = true

type exposer struct{ cmdt *commandTemplate }

func closure(cmdt commandTemplate, args ...interface{}) Shfn {
	if len(args) == 0 {
		// an empty call is a synonym for Shfn.Run().
		// if you want to just get a RunningCommand reference to track, use Shfn.Start() instead.
		enclose(&cmdt).Run()
		return nil
	} else if args[0] == expose {
		// produce a function that when called with an exposer, exposes its cmdt.
		return func(x ...interface{}) Shfn {
			t := x[0].(*exposer)
			t.cmdt = &cmdt
			return nil
		}
	} else {
		// examine each of the arguments, modify our (already forked) cmdt, and
		//  return a new callable Shfn closure with the newly baked command template.
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

func (f Shfn) expose() *commandTemplate {
	var t exposer
	f(expose)(&t)
	return t.cmdt
}

func enclose(cmdt *commandTemplate) Shfn {
	return func(x ...interface{}) Shfn {
		return closure(*cmdt, x...)
	}
}

func (f Shfn) BakeArgs(args ...string) Shfn {
	return enclose(f.expose().bakeArgs(args...))
}

func (cmdt *commandTemplate) bakeArgs(args ...string) *commandTemplate {
	cmdt.args = append(cmdt.args, args...)
	return cmdt
}

func (f Shfn) BakeEnv(args Env) Shfn {
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

func (f Shfn) ClearEnv() Shfn {
	return enclose(f.expose().clearEnv())
}

func (cmdt *commandTemplate) clearEnv() *commandTemplate {
	cmdt.env = make(map[string]string)
	return cmdt
}

func (f Shfn) BakeOpts(args ...Opts) Shfn {
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
		if arg.OkExit != nil {
			cmdt.OkExit = arg.OkExit
		}
	}
	return cmdt
}

/**
 * Starts execution of the command.  Returns a reference to a RunningCommand,
 * which can be used to track execution of the command, configure exit listeners,
 * etc.
 */
func (f Shfn) Start() *RunningCommand {
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
		case Shfn:
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

/**
 * Starts execution of the command, and waits until completion before returning.
 *
 * The is exactly the behavior of a no-arg invokation on an Shfn, i.e.
 *   `Sh("echo")()`
 * and
 *   `Sh("echo").Run()`
 * are interchangable and behave identically.
 *
 * Use the Start() method instead if you need to run a task in the background, or
 * if you otherwise need greater control over execution.
 */
func (f Shfn) Run() {
	cmdt := f.expose()
	cmd := f.Start()
	cmd.Wait()
	exitCode := cmd.GetExitCode()
	for _, okcode := range cmdt.OkExit {
		if exitCode == okcode {
			return
		}
	}
	panic(FailureExitCode{cmdname: cmdt.cmd, code: exitCode})
}
