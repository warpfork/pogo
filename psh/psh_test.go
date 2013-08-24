package psh

import (
	"os/exec"
	"testing"
	"time"
)

func TestStateConsts(t *testing.T) {
	if !(UNSTARTED == 0 && RUNNING == 1 && FINISHED == 2 && PANICED == 3) {
		t.Fail()
	}
}

func TestPshExecBasic(t *testing.T) {
	cmdr := NewRunningCommand(
		exec.Command("echo"),
	)
	cmdr.startCalmly()
	cmdr.WaitSoon(1 * time.Second)
	cmdr.GetExitCode()
}
