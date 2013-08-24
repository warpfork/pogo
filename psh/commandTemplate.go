package psh

import (
	"fmt"
	"os/exec"
)

func Command(cmd string) *commandTemplate {
	return &commandTemplate{
		cmd: cmd,
	}
}

func (cmdt commandTemplate) Bake(args ...string) *commandTemplate {
	cmdt.args = append(cmdt.args, args...)
	return &cmdt
}

func (cmdt commandTemplate) BakeEnv(env map[string]string) *commandTemplate {
	//TODO
	return &cmdt
}

func (cmdt commandTemplate) BakeOpts(opts Opts) *commandTemplate {
	//TODO
	return &cmdt
}

func (cmdt commandTemplate) Start() *RunningCommand {
	cmd := NewRunningCommand(exec.Command(cmdt.cmd, cmdt.args...))
	cmd.Start()
	return cmd
}

/**
 * Executes the command, waits for it, and panics with FailureExitCode if it exits with an unexpected code.
 */
func (cmdt commandTemplate) Go() {

// i think it's entirely possible that this should take args of (moar ...interface{}), which can be arg strings, or of type Opts (!) (env can go fuck itself, or at least maybe needs to be a type)

// you could just make all of the things return the same function again, and it always accepts whatever, and it's always baking until you call it with zero args.
// if you want it to return a RunningCommand instead of waiting, you set up `var RunningCommand{}` and your last call has a single arg of (&cmdr).
// come to think of it, you might want to have an UNINITIALIZED state for RunningCommand so we can expose it (for the reason on the line above) but never let people shoot themselves in the foot by waiting for something they initialized wrong.

	cmd := cmdt.Start()
	exitCode := cmd.GetExitCode()
	//TODO support configurable expected exit codes
	if exitCode != 0 {
		panic(FailureExitCode{cmdname: cmdt.cmd, code: exitCode})
	}
}

type Opts struct {
	Cwd string

	/**
	 * Can be a:
	 *   - string, in which case it will be copied in literally
	 *   - []byte, again, taken literally
	 *   - buffer, all that sort of thing, taken literally
	 *   - an io.Reader, in which case that will be streamed in
	 *   - another Command, in which case that wil be started with this one and its output piped into this one
	 */
	In interface{}
}


type commandTemplate struct {
	cmd string

	args []string

	env map[string]string

	Opts
}

type FailureExitCode struct {
	cmdname string
	code int
}

func (err FailureExitCode) Error() string {
	return fmt.Sprintf("command \"%s\" exited with unexpected status %d", err.cmdname, err.code)
}
