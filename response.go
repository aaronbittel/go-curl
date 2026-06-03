package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Response struct {
	StatusLine *StatusLine
	Header     http.Header
	Body       io.ReadCloser
}

type StatusLine struct {
	Version    string
	StatusCode int
	StatusText string
}

type ResponseBody struct {
	r    io.Reader
	conn net.Conn
}

func (body *ResponseBody) Read(p []byte) (n int, err error) {
	return body.r.Read(p)
}

func (body *ResponseBody) Close() error {
	return body.conn.Close()
}

func (s StatusLine) String() string {
	return fmt.Sprintf("%s %d %s", s.Version, s.StatusCode, s.StatusText)
}

func DumpResponse(resp *Response) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n", resp.StatusLine))
	sb.WriteString(fmt.Sprintf("Date: %s\n", resp.Header.Get("Date")))
	sb.WriteString(fmt.Sprintf("Content-Type: %s\n", resp.Header.Get("Content-Type")))
	sb.WriteString(fmt.Sprintf("Content-Length: %s\n", resp.Header.Get("Content-Length")))
	sb.WriteString(fmt.Sprintf("Connection: %s\n", resp.Header.Get("Connection")))
	sb.WriteString(fmt.Sprintf("Server: %s\n", resp.Header.Get("Server")))
	sb.WriteString(fmt.Sprintf("Access-Control-Allow-Origin: %s\n", resp.Header.Get("Access-Control-Allow-Origin")))
	sb.WriteString(fmt.Sprintf("Access-Control-Allow-Credentials: %s\n", resp.Header.Get("Access-Control-Allow-Credentials")))
	return sb.String()
}
