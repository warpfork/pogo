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

package iox

import (
	"github.com/coocood/assrt"
	"io"
	"sync"
	"testing"
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

	assert.Equal("asdf", <-ch)
	assert.Equal("", <-ch)
	assert.Equal("\nwakawaka", <-ch)
	assert.Equal("\tz", <-ch)
	_, open := <-ch
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

	assert.Equal("asdf", <-ch)
	_, open := <-ch
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

	assert.Equal([]byte("asdf"), <-ch)
	assert.Equal([]byte(""), <-ch)
	assert.Equal([]byte("\nwakawaka"), <-ch)
	assert.Equal([]byte("\tz"), <-ch)
	_, open := <-ch
	assert.Equal(false, open)
}
