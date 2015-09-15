package safe

import (
	"io"
	"sync"
)

// ChanWriter ...
type ChanWriter struct {
	C chan []byte
}

// NewChanWriter ...
func NewChanWriter(c chan []byte) *ChanWriter {
	s := ChanWriter{
		C: c,
	}

	return &s
}

// Write actually write the bytes to the writer
func (w *ChanWriter) Write(p []byte) (n int, err error) {
	w.C <- p

	return len(p), nil
}

// Writer allows safe writing on a writer
type Writer struct {
	Writer io.Writer
	Lock   sync.Mutex
}

// NewWriter ...
func NewWriter(w io.Writer) *Writer {
	s := Writer{
		Writer: w,
	}

	return &s
}

// Write actually write the bytes to the writer
func (w *Writer) Write(p []byte) (n int, err error) {
	w.Lock.Lock()
	defer w.Lock.Unlock()

	return w.Writer.Write(p)
}

// SeekWriter writes line one line at a time
type SeekWriter struct {
	R      rune
	BS     []byte
	Writer io.Writer
}

// NewSeekWriter ...
func NewSeekWriter(r rune, w io.Writer) *SeekWriter {
	s := SeekWriter{
		Writer: w,
		R:      r,
	}

	return &s
}

// Write actually write the bytes to the writer
func (w *SeekWriter) Write(p []byte) (n int, err error) {
	I := 0
	P := string(p)
	L := len(P)

	for i, b := range P {
		if b == w.R {
			line := append(w.BS, []byte(P[:i+1])...)
			n, err := w.Writer.Write(line)

			if err != nil {
				return n, err
			}

			w.BS = []byte(P[i+1:])
			I = i
		}
	}

	w.BS = append(w.BS, P[I:L]...)

	return L, nil
}
