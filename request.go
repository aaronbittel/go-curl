package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

type Method int

const Version = "HTTP/1.1"
const CRLF = "\r\n"

type Request struct {
	Method string
	Url    *Url
	Header http.Header
	Body   string
}

func NewRequest(method string, url *Url, body string) *Request {
	header := http.Header{}
	header.Set("Host", url.host)
	header.Set("Accept", "*/*")
	header.Set("Connection", "close")

	return &Request{Method: method, Url: url, Header: header, Body: body}
}

func (req Request) Send() (*Response, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", req.Url.host, req.Url.port))
	if err != nil {
		return nil, err
	}

	conn.Write([]byte(req.build()))

	br := bufio.NewReader(conn)
	statusLineRaw, err := readHttpLine(br)
	if err != nil {
		return nil, err
	}
	statusLine, err := parseStatusLine(statusLineRaw[:len(statusLineRaw)-2])
	if err != nil {
		return nil, fmt.Errorf("illegal status line %q: %s", statusLineRaw, err)
	}

	header := http.Header{}

	for {
		nextTwo, err := br.Peek(len(CRLF))
		if err != nil {
			return nil, err
		}
		if bytes.Equal(nextTwo, []byte(CRLF)) {
			discarded, err := br.Discard(len(CRLF))
			if err != nil {
				return nil, errors.New("discarding '\r\n' failed")
			}
			if discarded != len(CRLF) {
				return nil, errors.New("discarding '\r\n' failed")
			}
			break
		}
		headerLine, err := readHttpLine(br)
		if err != nil {
			return nil, err
		}

		key, value, found := strings.Cut(headerLine, ":")
		if !found {
			return nil, fmt.Errorf("illegal header %q, missing \":\"", headerLine)
		}
		header.Set(strings.TrimSpace(key), strings.TrimSpace(value))
	}

	var body io.ReadCloser = io.NopCloser(br)

	if hasChunked(header.Get("Transfer-Encoding")) {
		body = io.NopCloser(httputil.NewChunkedReader(body))
	}

	if encoding := header.Get("Content-Encoding"); encoding != "" {
		switch encoding {
		case "gzip":
			body, err = gzip.NewReader(body)
			if err != nil {
				return nil, fmt.Errorf("gzip decoding failed: %w", err)
			}
		case "bzip2":
			body = io.NopCloser(bzip2.NewReader(body))
		case "zlib":
			body, err = zlib.NewReader(body)
			if err != nil {
				return nil, fmt.Errorf("zlib decoding failed: %w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported compression method %q", encoding)
		}
	}

	return &Response{
		StatusLine: statusLine,
		Header:     header,
		Body:       body,
	}, nil
}

func (req Request) RequestString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("connecting to %s\n", req.Url.host))
	sb.WriteString(fmt.Sprintf("Sending request %s /%s HTTP/1.1\n", req.Method, req.Url.path))
	for k, vs := range req.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(vs, ", ")))
	}
	return sb.String()
}

func (req Request) build() string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("%s /%s %s\r\n", req.Method, req.Url.path, Version))
	for k, vs := range req.Header {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(vs, ", ")))
	}

	if req.Body != "" {
		buf.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(req.Body)))
	}

	buf.WriteString("\r\n")
	buf.WriteString(req.Body)

	return buf.String()
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

func hasChunked(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "chunked")
}
