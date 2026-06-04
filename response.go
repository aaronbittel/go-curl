package main

import (
	"fmt"
	"io"
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

func (s StatusLine) String() string {
	return fmt.Sprintf("%s %d %s", s.Version, s.StatusCode, s.StatusText)
}

func DumpResponse(resp *Response) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n", resp.StatusLine))
	for k, vs := range resp.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(vs, ", ")))
	}
	return sb.String()
}
