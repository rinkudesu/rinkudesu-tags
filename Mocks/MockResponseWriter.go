package Mocks

import (
	"bufio"
	"net"
	"net/http"
)

type NilResponseWriter struct {
}

func (n NilResponseWriter) Header() http.Header {
	return nil
}

func (n NilResponseWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (n NilResponseWriter) WriteHeader(statusCode int) {
}

func (n NilResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	panic("not supported")
}

func (n NilResponseWriter) Flush() {
}

func (n NilResponseWriter) CloseNotify() <-chan bool {
	panic("not supported")
}

func (n NilResponseWriter) Status() int {
	return 0
}

func (n NilResponseWriter) Size() int {
	return 0
}

func (n NilResponseWriter) WriteString(s string) (int, error) {
	return 0, nil
}

func (n NilResponseWriter) Written() bool {
	return false
}

func (n NilResponseWriter) WriteHeaderNow() {
}

func (n NilResponseWriter) Pusher() http.Pusher {
	panic("not supported")
}
