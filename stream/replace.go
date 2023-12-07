// Package stream provides a set of utils to modify the stream.
package stream

import (
	"bytes"
	"io"
)

// NewSimpleReplacer returns a new Reader with streaming replace supports.
// It will replace every sequences matching `old` with `new`.
// This replacer will try to use as less mem as possible.
func NewSimpleReplacer(r io.Reader, old []byte, new []byte) io.Reader {
	return &simpleReplacer{
		underlying: r,

		old:    old,
		new:    new,
		oldLen: len(old),
	}
}

type simpleReplacer struct {
	underlying io.Reader

	old, new []byte
	oldLen   int

	buf        []byte
	bufLen     int
	overlapped int
	lastErr    error
}

func (r *simpleReplacer) consumeBuffer(p []byte) (n int, err error) {
	if r.lastErr != nil {
		r.overlapped = 0
	}
	if r.bufLen > r.overlapped {
		n = copy(p, r.buf[:r.bufLen-r.overlapped])
		r.bufLen = copy(r.buf, r.buf[n:r.bufLen])
	}
	if r.bufLen == 0 {
		err = r.lastErr
	}
	return
}

func (r *simpleReplacer) Read(p []byte) (n int, err error) {
	if n, err = r.consumeBuffer(p); n > 0 || err != nil {
		return
	}

	// fast path
	n, err = r.underlying.Read(p)
	if r.oldLen == 0 || (r.bufLen == 0 && bytes.IndexByte(p[:n], r.old[0]) == -1) {
		return
	}

	// slow path
	r.buf = append(r.buf[:r.bufLen], p[:n]...)
	cur := 0
	for {
		idx := bytes.Index(r.buf[cur:], r.old)
		if idx < 0 {
			break
		}

		cur += idx + r.oldLen
	}

	r.overlapped = overlapLen(r.buf[cur:], r.old, r.oldLen-1)
	if cur > 0 {
		// copy happening
		r.buf = bytes.ReplaceAll(r.buf, r.old, r.new)
	}

	r.bufLen = len(r.buf)

	r.lastErr = err

	return r.consumeBuffer(p)
}

// overlapLen returns how many characters being overlapped.
// For instance, if we have the following 2 sequences:
//
//	l: [a|b|c|d|e]
//	r:     [c|d|e|f|g]
//
// then the overlap length would be 3
func overlapLen(l, r []byte, maxLen int) int {
	for ; maxLen > 0; maxLen-- {
		if maxLen > len(l) || maxLen > len(r) {
			continue
		}

		if bytes.Equal(l[len(l)-maxLen:], r[:maxLen]) {
			return maxLen
		}
	}

	return 0
}
