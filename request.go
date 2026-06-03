package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Method int

const Version = "HTTP/1.1"

type Request struct {
	Method string
	Url    *Url
	Header http.Header
}

func (req Request) build() string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s /%s %s\r\n", req.Method, req.Url.path, Version))
	for k, vs := range req.Header {
		for _, v := range vs {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	buf.WriteString("\r\n")

	return buf.String()
}

func (req Request) RequestString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("connecting to %s\n", req.Url.host))
	sb.WriteString(fmt.Sprintf("Sending request %s /%s HTTP/1.1\n", req.Method, req.Url.path))
	for k, vs := range req.Header {
		for _, v := range vs {
			sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}
	return sb.String()
}
