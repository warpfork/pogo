package iox

import (
"fmt"
	"bytes"
	"io"
)

func WriterFromInterface(x interface{}) io.Writer {
	switch y := x.(type) {
	case []byte:
		return WriterToByteSlice(y)
	case io.Writer:
		return y
	case bytes.Buffer:
		return &y
	case chan<- string:
		return WriterToChanString(y)
	case chan<- []byte:
		return WriterToChanByteSlice(y)
	default:
		return nil
	}
}

func WriterToByteSlice(bats []byte) io.Writer {
	return &writerByteSlice{slice: bats}
}

type writerByteSlice struct {
	slice []byte
}

func (r *writerByteSlice) Write(p []byte) (n int, err error) {
	c := cap(r.slice)
	w := len(r.slice)
	available := c - w
	if len(p) < available {
		available = len(p)
	}
	r.slice = r.slice[:w+available]
	fmt.Printf(":: %d %d %d \n", w, c, available)
	n = copy(r.slice[w:], p[:available])
	fmt.Printf("::: %d \n", n)
	return
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
