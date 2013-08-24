package psh

import (
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	/**
	 * 'Unstarted' is the state of a command that has been constructed, but execution has not yet begun.
	 */
	UNSTARTED int32 = iota

	/**
	 * 'Running' is the state of a command that has begun execution, but not yet finished.
	 */
	RUNNING

	/**
	 * 'Finished' is the state of a command that has finished normally.
	 *
	 * The exit code may or may not have been success, but at the very least we
	 * successfully observed that exit code.
	 */
	FINISHED

	/**
	 * 'Panicked' is the state of a command that at some point began execution, but has encountered
	 * serious problems.
	 *
	 * It may not be clear whether or not the command is still running, since a panic implies we no
	 * longer have completely clear visibility to the command on the underlying system.  The exit
	 * code may not be reliably known.
	 */
	PANICKED
)

func NewRunningCommand(cmd *exec.Cmd) *RunningCommand {
	return &RunningCommand{
		cmd:      cmd,
		state:    UNSTARTED,
		exitCh:   make(chan bool),
		exitCode: -1,
	}
}

type RunningCommand struct {
	mutex sync.Mutex

	/** Always access this with functions from the atomic package, and when
	* transitioning states set the status after all other fields are mutated,
	* so that checks of State() serve as a memory barrier for all. */
	state int32

	cmd *exec.Cmd
	//TODO: or: callback func() int

	/** If this is set, game over. */
	err error

	/** Wait for this to close in order to wait for the process to return. */
	exitCh chan bool

	/** Exit code if we're state==FINISHED and exit codes are possible on this platform, or
	 * -1 if we're not there yet.  Will not change after exitCh has closed. */
	exitCode int

	/** Functions to call back when the command has exited. */
	exitListeners []func(*RunningCommand)
}

func (cmd *RunningCommand) State() int32 {
	return atomic.LoadInt32(&cmd.state)
}

/** Returns true if the command is current running. */
func (cmd *RunningCommand) IsRunning() bool {
	state := cmd.State()
	return state == RUNNING
}

/** Returns true if the command has ever been started (including if the command is already finished). */
func (cmd *RunningCommand) IsStarted() bool {
	state := cmd.State()
	return state == RUNNING || state == FINISHED || state == PANICKED
}

/** Returns true if the command is finished (either gracefully, or with internal errors). */
func (cmd *RunningCommand) IsDone() bool {
	state := cmd.State()
	return state == FINISHED || state == PANICKED
}

/** Returns true if the command is finished either gracefully.  (A nonzero exit code may still be set.) */
func (cmd *RunningCommand) IsFinishedGracefully() bool {
	state := cmd.State()
	return state == FINISHED
}

func (cmd *RunningCommand) Start() {
	if err := cmd.startCalmly(); err != nil {
		panic(err)
	}
}

func (cmd *RunningCommand) startCalmly() error {
	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if cmd.IsStarted() {
		return nil
	}

	atomic.StoreInt32(&cmd.state, RUNNING)
	if err := cmd.cmd.Start(); err != nil {
		cmd.finalState(CommandStartError{cause: err})
		return cmd.err
	}

	go cmd.waitAndHandleExit()
	return nil
}

func (cmd *RunningCommand) waitAndHandleExit() {
	err := cmd.cmd.Wait()

	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if err == nil {
		cmd.exitCode = 0
	} else if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			cmd.exitCode = waitStatus.ExitStatus()
		} else {
			panic(exitError) //TODO: damage control better.  consider setting some kind of CommandMonitorError.
		}
	} else {
		panic(err) //TODO: damage control better.  consider setting some kind of CommandMonitorError.
	}
	cmd.finalState(nil)
}

func (cmd *RunningCommand) finalState(err error) {
	// must hold cmd.mutex before calling this
	// golang is an epic troll: claims to be best buddy for concurrent code, SYNC PACKAGE DOES NOT HAVE REENTRANT LOCKS
	if cmd.IsRunning() {
		if err == nil {
			atomic.StoreInt32(&cmd.state, FINISHED)
		} else {
			cmd.err = err
			atomic.StoreInt32(&cmd.state, PANICKED)
		}
		//TODO iterate over exit listeners
		for _, cb := range cmd.exitListeners {
			func() {
				defer recover()
				cb(cmd)
			}()
		}
	}
	close(cmd.exitCh)
}

/**
 * Add a function to be called when this command completes.
 *
 * These listener functions will be invoked after the exit code and other command
 * state is final, but before other Wait() methods unblock.
 * (This means if you want for example to log a message that a process exited, and
 * your main function is Wait()'ing for that process... if you use AddExitListener()
 * to invoke your log function then you will always get the log.)
 *
 * The listener function should complete quickly and not try to perform other blocking
 * operations or locks, since other actions are waiting until the listeners have all
 * been called.  Panics that escape the function will be silently discarded; do not
 * panic in a listener.
 *
 * If the command is already in the state FINISHED or PANICKED, the callback function
 * will be invoked immediately in the current goroutine.
 */
func (cmd *RunningCommand) AddExitListener(callback func(*RunningCommand)) {
	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if cmd.IsDone() {
		func() {
			defer recover()
			callback(cmd)
		}()
	} else {
		cmd.exitListeners = append(cmd.exitListeners, callback)
	}
}

/**
 * Returns a channel that will be open until the command is complete.
 * This is suitable for use in a select block.
 */
func (cmd *RunningCommand) GetExitChannel() <-chan bool {
	return cmd.exitCh
}

/**
 * Waits for the command to exit before returning.
 *
 * There are no consequences to waiting on a single command repeatedly;
 * all wait calls will return normally when the command completes.  The order
 * in which multiple wait calls will return is undefined.  Similarly, there
 * are no consequences to waiting on a command that has not yet started;
 * the function will still wait without error until the command finishes.
 * (Much friendlier than os.exec.Cmd.Wait(), neh?)
 */
func (cmd *RunningCommand) Wait() {
	<-cmd.GetExitChannel()
}

/**
 * Waits for the command to exit before returning, or for the specified duration.
 * Returns true if the return was due to the command finishing, or false if the
 * return was due to timeout.
 */
func (cmd *RunningCommand) WaitSoon(d time.Duration) bool {
	select {
	case <-time.After(d):
		return false
	case <-cmd.GetExitChannel():
		return true
	}
}

/**
 * Waits for the command to exit if it has not already, then returns the exit code.
 */
func (cmd *RunningCommand) GetExitCode() int {
	if !cmd.IsDone() {
		cmd.Wait()
	}
	return cmd.exitCode
}

/**
 * Waits for the command to exit if it has not already, or for the specified duration,
 * then either returns the exit code, or -1 if the duration expired and the command
 * still hasn't returned.
 */
func (cmd *RunningCommand) GetExitCodeSoon(d time.Duration) int {
	if cmd.WaitSoon(d) {
		return cmd.exitCode
	} else {
		return -1
	}
}
