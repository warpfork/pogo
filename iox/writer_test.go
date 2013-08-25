package iox

import (
	"github.com/coocood/assrt"
	"io"
	"testing"
	"sync"
)

func TestWriterToChanString(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan string)
	w := WriterToChanString(ch)
	go func() {
		w.Write([]byte("asdf"))
		w.Write([]byte(""))
		w.Write([]byte("\nwakawaka"))
		w.Write([]byte("\tz"))
		close(ch)
	}()

	assert.Equal("asdf", <- ch)
	assert.Equal("", <- ch)
	assert.Equal("\nwakawaka", <- ch)
	assert.Equal("\tz", <- ch)
	_, open := <- ch
	assert.Equal(false, open)
}

func TestWriterToChanStringClosed(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan string)
	w := WriterToChanString(ch)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		n, err := w.Write([]byte("asdf"))
		assert.Equal(4, n)
		assert.Equal(nil, err)
		close(ch)
		n, err = w.Write([]byte("\tz"))
		assert.Equal(0, n)
		assert.Equal(io.EOF, err)
		wg.Done()
	}()

	assert.Equal("asdf", <- ch)
	_, open := <- ch
	assert.Equal(false, open)
	wg.Wait()
}

func TestWriterToChanByteSlice(t *testing.T) {
	assert := assrt.NewAssert(t)

	ch := make(chan []byte)
	w := WriterToChanByteSlice(ch)
	go func() {
		w.Write([]byte("asdf"))
		w.Write([]byte(""))
		w.Write([]byte("\nwakawaka"))
		w.Write([]byte("\tz"))
		close(ch)
	}()

	assert.Equal([]byte("asdf"), <- ch)
	assert.Equal([]byte(""), <- ch)
	assert.Equal([]byte("\nwakawaka"), <- ch)
	assert.Equal([]byte("\tz"), <- ch)
	_, open := <- ch
	assert.Equal(false, open)
}
