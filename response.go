package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
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

func NewRequest(method string, url *Url) (*Request, error) {
	header := http.Header{}
	header.Set("Host", url.host)
	header.Set("Accept", "*/*")
	header.Set("Connection", "close")

	m, err := ParseMethod(method)
	if err != nil {
		return nil, err
	}

	return &Request{Method: m, Url: url, Header: header}, nil
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
	statusLine, err := ParseStatusLine(statusLineRaw[:len(statusLineRaw)-2])
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

func (body *ResponseBody) Read(p []byte) (n int, err error) {
	return body.r.Read(p)
}

func (body *ResponseBody) Close() error {
	return body.conn.Close()
}

func ParseStatusLine(statusLine string) (*StatusLine, error) {
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
