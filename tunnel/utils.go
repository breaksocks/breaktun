package tunnel

import (
	"errors"
	"io"
	"net"
	"time"
)

func WriteN2(bs []byte, offset int, n uint16) {
	bs[offset] = byte(n >> 8)
	bs[offset+1] = byte(n & 0xFF)
}

func ReadN2(bs []byte, offset int) uint16 {
	return (uint16(bs[offset]) << 8) | uint16(bs[offset+1])
}

func WriteN4(bs []byte, offset int, n uint32) {
	bs[offset] = byte(n >> 24)
	bs[offset+1] = byte(n >> 16)
	bs[offset+2] = byte(n >> 8)
	bs[offset+3] = byte(n)
}

func ReadN4(bs []byte, offset int) uint32 {
	var n uint32
	n |= uint32(bs[offset]) << 24
	n |= uint32(bs[offset+1]) << 16
	n |= uint32(bs[offset+2]) << 8
	n |= uint32(bs[offset+3])
	return n
}

type RetryRequester struct {
	read, write chan []byte
	maxRetry    int
	timeout     int
}

func NewRetryRequester(max_retry, timeout int, read, write chan []byte) *RetryRequester {
	return &RetryRequester{
		read: read, write: write,
		maxRetry: max_retry, timeout: timeout,
	}
}

var RequestTimeout error = errors.New("request timeout")

func (req *RetryRequester) Request(data []byte) ([]byte, error) {
	for i := 0; i < req.maxRetry; i += 1 {
		req.write <- data
		select {
		case rep, ok := <-req.read:
			if ok {
				return rep, nil
			}
			return nil, io.EOF
		case <-time.After(time.Millisecond * req.timeout):
		}
	}
	return nil, RequestTimeout
}

func (req *RetryRequester) RetryGet(max_retry int) ([]byte, error) {
	for i := 0; i < max_retry; i += 1 {
		select {
		case rep, ok := <-req.read:
			if ok {
				return rep, nil
			}
			return nil, io.EOF
		case <-time.After(time.Millisecond * req.timeout):
		}
	}
	return nil, RequestTimeout
}
