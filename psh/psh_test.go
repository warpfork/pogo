package psh

import (
	"github.com/coocood/assrt"
	"os/exec"
	"testing"
	"time"
)

func TestStateConsts(t *testing.T) {
	if !(UNSTARTED == 0 && RUNNING == 1 && FINISHED == 2 && PANICED == 3) {
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
}
