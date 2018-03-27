package feiniubus

import (
	"io"
	"sync"
)

type offsetReader struct {
	buf    io.ReadSeeker
	lock   sync.Mutex
	closed bool
}

func newOffsetReader(buf io.ReadSeeker, offset int64) *offsetReader {
	reader := &offsetReader{}
	buf.Seek(offset, 0)

	reader.buf = buf
	return reader
}

func (o *offsetReader) Close() error {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.closed = true
	return nil
}

func (o *offsetReader) Read(p []byte) (int, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.closed {
		return 0, io.EOF
	}

	return o.buf.Read(p)
}

func (o *offsetReader) Seek(offset int64, whence int) (int64, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	return o.buf.Seek(offset, whence)
}

func (o *offsetReader) CloseAndCopy(offset int64) *offsetReader {
	o.Close()
	return newOffsetReader(o.buf, offset)
}
