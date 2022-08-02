package Mocks

import (
	"bufio"
	"net"
	"net/http"
)

const (
	nilResponseWriterNotSupportedMessage = "not supported"
)

type NilResponseWriter struct {
}

func (n NilResponseWriter) Header() http.Header {
	return nil
}

func (n NilResponseWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

func (n NilResponseWriter) WriteHeader(_ int) {
	// do nothing
}

func (n NilResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	panic(nilResponseWriterNotSupportedMessage)
}

func (n NilResponseWriter) Flush() {
	// do nothing
}

func (n NilResponseWriter) CloseNotify() <-chan bool {
	panic(nilResponseWriterNotSupportedMessage)
}

func (n NilResponseWriter) Status() int {
	return 0
}

func (n NilResponseWriter) Size() int {
	return 0
}

func (n NilResponseWriter) WriteString(_ string) (int, error) {
	return 0, nil
}

func (n NilResponseWriter) Written() bool {
	return false
}

func (n NilResponseWriter) WriteHeaderNow() {
	// do nothing
}

func (n NilResponseWriter) Pusher() http.Pusher {
	panic(nilResponseWriterNotSupportedMessage)
}
