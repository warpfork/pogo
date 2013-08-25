package psh

type commandTemplate struct {
	cmd string

	args []string

	env Env

	Opts
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

type Env map[string]string

type ClearEnv struct{}
