package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Method int

const Version = "HTTP/1.1"

type Request struct {
	Method string
	Url    *Url
	Header http.Header
}

func NewRequest(method string, url *Url) *Request {
	header := http.Header{}
	header.Set("Host", url.host)
	header.Set("Accept", "*/*")
	header.Set("Connection", "close")

	return &Request{Method: method, Url: url, Header: header}
}

func (req Request) Send() (*Response, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", req.Url.host, req.Url.port))
	if err != nil {
		return nil, err
	}

	conn.Write([]byte(req.build()))

	r := bufio.NewReader(conn)
	statusLineRaw, err := readHttpLine(r)
	if err != nil {
		return nil, err
	}
	statusLine, err := parseStatusLine(statusLineRaw[:len(statusLineRaw)-2])
	if err != nil {
		return nil, fmt.Errorf("illegal status line %q: %s", statusLineRaw, err)
	}

	header := http.Header{}

	for {
		nextTwo, err := r.Peek(2)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(nextTwo, []byte("\r\n")) {
			break
		}
		headerLine, err := readHttpLine(r)
		if err != nil {
			return nil, err
		}

		key, value, found := strings.Cut(headerLine, ":")
		if !found {
			return nil, fmt.Errorf("illegal header %q, missing \":\"", headerLine)
		}
		header.Set(strings.TrimSpace(key), strings.TrimSpace(value))
	}

	return &Response{
		StatusLine: statusLine,
		Header:     header,
		Body: &ResponseBody{
			r:    r,
			conn: conn,
		},
	}, nil
}

func parseStatusLine(statusLine string) (*StatusLine, error) {
	parts := strings.SplitN(statusLine, " ", 3)

	if len(parts) != 3 {
		return nil, fmt.Errorf("status line must contain exactly 3 parts")
	}

	version := parts[0]

	if version != Version {
		return nil, fmt.Errorf("illegal http version %q, expected %q", version, Version)
	}

	statusCode, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("illegal status code: %q", parts[1])
	}

	statusText := http.StatusText(statusCode)
	if statusText == "" {
		return nil, fmt.Errorf("unknown status code: %q", parts[1])
	}

	if !strings.EqualFold(statusText, parts[2]) {
		return nil, fmt.Errorf("illegal status text: %q, expected %q", parts[2], statusText)
	}

	return &StatusLine{
		Version:    version,
		StatusCode: statusCode,
		StatusText: statusText,
	}, nil
}

func readHttpLine(r *bufio.Reader) (string, error) {
	httpLine, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(httpLine, "\r\n") {
		return "", fmt.Errorf("missing \"\r\n\"")
	}
	return httpLine, nil
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
