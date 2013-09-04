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
	"bytes"
	"io"
)

func WriterFromInterface(x interface{}) io.Writer {
	switch y := x.(type) {
	case io.Writer:
		return y
	case bytes.Buffer:
		return &y
	case chan<- string:
		return WriterToChanString(y)
	case chan string:
		return WriterToChanString(y)
	case chan<- []byte:
		return WriterToChanByteSlice(y)
	case chan []byte:
		return WriterToChanByteSlice(y)
	default:
		return nil
	}
}

func WriterToChanString(ch chan<- string) io.Writer {
	return &writerChanString{ch: ch}
}

type writerChanString struct {
	ch chan<- string
}

func (r *writerChanString) Write(p []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			n = 0
			err = io.EOF
		}
	}()

	r.ch <- string(p)
	return len(p), nil
}

func WriterToChanByteSlice(ch chan<- []byte) io.Writer {
	return &writerChanByteSlice{ch: ch}
}

type writerChanByteSlice struct {
	ch chan<- []byte
}

func (r *writerChanByteSlice) Write(p []byte) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			n = 0
			err = io.EOF
		}
	}()

	r.ch <- p
	return len(p), nil
}
