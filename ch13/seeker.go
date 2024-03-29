package ch13

import (
	"bytes"
	"errors"
	"io"
)

const (
	DefaultMaxBufferSize int64 = 1024 * 1024
)

var (
	ErrInvalidWhence    = errors.New("ch13.Seek: invalid whence")
	ErrNegativePosition = errors.New("ch13.Seek: negative position")
	ErrMaxSizeExceeded  = errors.New("ch13.Seek: max size exceeded")
)

type readSeeker struct {
	r           io.Reader
	buf         *bytes.Buffer
	i           int64
	max         int64
	completed   bool
	bufferError bool
}

func (r readSeeker) availableBuffer() int64 {
	return r.max - int64(r.buf.Len())
}

func (r *readSeeker) Read(p []byte) (n int, err error) {
	if r.bufferError {
		return 0, ErrMaxSizeExceeded
	}
	needToRead := int64(len(p))
	i := 0
	if r.i < int64(r.buf.Len()) {
		rawBuf := r.buf.Bytes()
		sizeToRead := copy(p, rawBuf[r.i:])
		i += sizeToRead
		r.i += int64(sizeToRead)
		needToRead -= int64(sizeToRead)
	} else if int64(r.buf.Len()) < r.i {
		needToSkip := r.i - int64(r.buf.Len())
		size, err := io.CopyN(r.buf, r.r, needToSkip)
		if err == io.EOF || size < needToSkip {
			r.completed = true
			return 0, io.EOF
		}
	}
	if needToRead > 0 {
		var read int
		readSize := int64(len(p) - i)
		if readSize > r.availableBuffer() {
			return 0, ErrMaxSizeExceeded
		}
		read, err = r.r.Read(p[i : i+int(readSize)])
		if err == io.EOF {
			r.completed = true
		} else if r.availableBuffer() == 0 && read < len(p)-i {
			return 0, ErrMaxSizeExceeded
		}
		if int64(read) < needToRead {
			r.completed = true
		}
		r.buf.Write(p[i : i+read])
		r.i += int64(read)
		i += read
	}
	return i, err
}

func (r *readSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
		r.bufferError = false
	case io.SeekCurrent:
		abs = r.i + offset
		r.bufferError = false
	case io.SeekEnd:
		if !r.completed {
			io.CopyN(r.buf, r.r, r.max-int64(r.buf.Len()))
			if r.availableBuffer() == 0 {
				r.bufferError = true
			}
			r.completed = true
		}
		abs = int64(r.buf.Len()) + offset
	default:
		return 0, ErrInvalidWhence
	}
	if abs < 0 {
		return 0, ErrNegativePosition
	}
	r.i = abs
	return abs, nil
}

type Option func(*readSeeker)

func MaxBufferSize(s uint64) Option {
	return func(rs *readSeeker) {
		rs.max = int64(s)
	}
}

func NewSeeker(r io.Reader, options ...Option) io.ReadSeeker {
	rs := &readSeeker{
		r:           r,
		buf:         &bytes.Buffer{},
		i:           0,
		max:         DefaultMaxBufferSize,
		completed:   false,
		bufferError: false,
	}
	for _, opt := range options {
		opt(rs)
	}
	return rs
}
