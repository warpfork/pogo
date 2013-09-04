// Copyright 2013 Eric Myhre
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// the go race detector would flag this as a race.
	// and correctly so!  that's why this field is private.
	// assert.Equal(
	// 	-1,
	// 	cmdr.exitCode,
	// )
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

func TestPshExitCode(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmdr := NewRunningCommand(
		exec.Command("sh", []string{"-c", "exit 22"}...),
	)
	cmdr.startCalmly()
	cmdr.WaitSoon(1 * time.Second)
	assert.Equal(
		22,
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

func TestPshNonexistentCommandPanics(t *testing.T) {
	assert := assrt.NewAssert(t)

	cmdr := NewRunningCommand(
		exec.Command("/thishadbetternotbeacommand"),
	)
	cmdr.startCalmly()
	cmdr.WaitSoon(1 * time.Second)
	assert.Equal(
		-1,
		cmdr.GetExitCode(),
	)
	assert.NotEqual(
		nil,
		cmdr.err,
	)
	assert.Equal(
		PANICKED,
		cmdr.State(),
	)
}
