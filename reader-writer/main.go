package main

import (
	"fmt"
	"io"
	"strings"
)

type myReader struct {
	source     string
	sourceByte []byte
	coursor    int
}

func (r *myReader) Read(p []byte) (n int, err error) {
	n = 0
	for n < len(p) {
		if r.coursor >= len(r.sourceByte) {
			return n, io.EOF
		}
		p[n] = r.sourceByte[r.coursor]
		n++
		r.coursor++
	}
	return n, err
}

type myWriter struct {
	buffer []byte
}

func (w *myWriter) Write(p []byte) (n int, err error) {
	w.buffer = p
	return len(p), nil
}

func (w *myWriter) ToUpper() {
	fmt.Printf("Wrting %d bytes: %s \n", len(w.buffer), strings.ToUpper(string(w.buffer)))
}

func main() {

	var source string

	source = "This is a long string that I want to read."
	sourceByte := []byte(source)
	destByte := make([]byte, 5)

	myNewReader := myReader{}
	myNewWriter := myWriter{}

	myNewReader.sourceByte = sourceByte
	var err error
	var n int
	for err == nil {
		n, err = myNewReader.Read(destByte)
		myNewWriter.Write(destByte[:n])
		myNewWriter.ToUpper()
	}

}
