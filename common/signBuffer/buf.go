package signBuffer

import (
	"bytes"
	"time"
)

type semaphore struct{}

type SignBuffer struct {
	sign chan semaphore
	*bytes.Buffer
	closeHook func(b []byte) bool
	// StreamTimeout
	readTimeout time.Duration
	// 过滤词
	filterWords []byte
}

func New(closeHook func(b []byte) bool, readTimeout time.Duration, filterWords []byte) *SignBuffer {
	return &SignBuffer{
		sign:        make(chan semaphore, 32),
		Buffer:      bytes.NewBuffer(nil),
		closeHook:   closeHook,
		readTimeout: readTimeout,
		filterWords: filterWords,
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
	n, err = s.Buffer.Write(bytes.ReplaceAll(p, s.filterWords, []byte{}))
	// 如果 p 含有 [DONE] 则关闭通道
	if s.closeHook(p) {
		s.Close()
	} else {
		s.sign <- semaphore{}
	}
	return
}

func (s *SignBuffer) Close() error {
	close(s.sign)
	return nil
}
