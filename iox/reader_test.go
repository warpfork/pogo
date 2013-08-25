package iox

import (
	"bytes"
	"github.com/coocood/assrt"
	"io"
	"testing"
)

func TestReaderFromChanString(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan string)
	var output bytes.Buffer
	go func() {
		ch <- "asdf"
		ch <- ""
		ch <- "\nwakawaka"
		ch <- "\tz"
		close(ch)
	}()
	io.Copy(&output, ReaderFromChanString(ch))

	assert.Equal(
		"asdf\nwakawaka\tz",
		output.String(),
	)
}

func TestReaderFromByteSlice(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan []byte)
	var output bytes.Buffer
	go func() {
		ch <- []byte("asdf")
		ch <- []byte("")
		ch <- []byte("\nwakawaka")
		ch <- []byte("\tz")
		close(ch)
	}()
	io.Copy(&output, ReaderFromChanByteSlice(ch))

	assert.Equal(
		"asdf\nwakawaka\tz",
		output.String(),
	)
}
