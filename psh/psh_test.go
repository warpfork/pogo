package psh

import (
	"github.com/coocood/assrt"
	"os/exec"
	"sync"
	"testing"
	"time"
)

func TestStateConsts(t *testing.T) {
	if !(UNSTARTED == 0 && RUNNING == 1 && FINISHED == 2 && PANICKED == 3) {
		t.Fail()
	}
}

// Test that we can exec something, wait, and it returns quickly and with an exit code.
func TestPshExecBasic(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmdr := NewRunningCommand(
		exec.Command("echo"),
	)
	cmdr.startCalmly()
	cmdr.WaitSoon(1 * time.Second)
	assert.Equal(
		0,
		cmdr.GetExitCode(),
	)
	assert.Equal(
		nil,
		cmdr.err,
	)
	assert.Equal(
		FINISHED,
		cmdr.State(),
	)
}

func TestPshWaitTimeout(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmdr := NewRunningCommand(
		exec.Command("sleep", "1"),
	)
	cmdr.startCalmly()
	assert.Equal(
		false,
		cmdr.WaitSoon(20*time.Millisecond),
	)
	assert.Equal(
		-1,
		cmdr.exitCode,
	)
	assert.Equal(
		RUNNING,
		cmdr.State(),
	)
}

func TestPshExitListeners(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmdr := NewRunningCommand(
		exec.Command("echo"),
	)
	cmdr.startCalmly()
	var wg sync.WaitGroup
	wg.Add(1)
	cmdr.AddExitListener(func(*RunningCommand) {
		defer wg.Done()
		assert.Equal(
			0,
			cmdr.GetExitCode(),
		)
		assert.Equal(
			nil,
			cmdr.err,
		)
		assert.Equal(
			FINISHED,
			cmdr.State(),
		)
	})
	wg.Wait()
	assert.Equal(
		0,
		cmdr.GetExitCode(),
	)
	assert.Equal(
		nil,
		cmdr.err,
	)
	assert.Equal(
		FINISHED,
		cmdr.State(),
	)
}
