package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Method int

const (
	GET Method = iota
)

const Version = "HTTP/1.1"

type Request struct {
	Method Method
	Url    *Url
	Header http.Header
}

func ParseMethod(method string) (Method, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return GET, nil
	default:
		return 0, fmt.Errorf("unknown method %q", method)
	}
}

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	default:
		return ""
	}
}

func (req Request) build() string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s /%s %s\r\n", req.Method, req.Url.path, Version))
	for k, v := range req.Header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	return buf.String()
}

func (req Request) RequestString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("connecting to %s\n", req.Url.host))
	sb.WriteString(fmt.Sprintf("Sending request GET /%s HTTP/1.1\n", req.Url.path))
	for k, v := range req.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	return sb.String()
}
