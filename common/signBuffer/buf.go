package signBuffer

import (
	"bytes"
	"time"
)

type SignBuffer struct {
	sign chan struct{}
	*bytes.Buffer
	closeHook func(b []byte) bool
	// StreamTimeout
	readTimeout time.Duration
}

func New(closeHook func(b []byte) bool, readTimeout time.Duration) *SignBuffer {
	return &SignBuffer{
		sign:        make(chan struct{}, 32),
		Buffer:      bytes.NewBuffer(nil),
		closeHook:   closeHook,
		readTimeout: readTimeout,
	}
}

func (s *SignBuffer) Read(p []byte) (n int, err error) {
	select {
	case <-s.sign:
	case <-time.After(s.readTimeout):
	}
	return s.Buffer.Read(p)
}

func (s *SignBuffer) Write(p []byte) (n int, err error) {
	n, err = s.Buffer.Write(bytes.ReplaceAll(p, []byte("data:"), []byte{}))
	// 如果 p 含有 [DONE] 则关闭通道
	if s.closeHook(p) {
		s.Close()
	} else {
		s.sign <- struct{}{}
	}
	return
}

func (s *SignBuffer) Close() error {
	close(s.sign)
	return nil
}
