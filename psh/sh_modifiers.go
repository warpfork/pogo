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
	 *   - io.Reader, which will be streamed in
	 *   - bytes.Buffer, all that sort of thing, taken literally
	 *   - <-chan string, in which case that will be streamed in
	 *   - <-chan byte[], in which case that will be streamed in
	 *   - another Command, in which case that will be started with this one and its output piped into this one
	 */
	In interface{}

	/**
	 * Can be a:
	 *   - []byte, which will be written to literally
	 *   - bytes.Buffer, which will be written to literally
	 *   - io.Writer, which will be written to streamingly, flushed to whenever the command flushes
	 *   - chan<- string, which will be written to streamingly, flushed to whenever a line break occurs in the output
	 *   - chan<- byte[], which will be written to streamingly, flushed to whenever the command flushes
	 *
	 * (There's nothing that's quite the equivalent of how you can give In a string, sadly; since
	 * strings are immutable in golang, you can't set Out=&str and get anywhere.)
	 */
	Out interface{}
}

type Env map[string]string

type ClearEnv struct{}
