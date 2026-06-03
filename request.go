package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Request struct {
	Req *http.Request
	Url *Url
}

func NewRequest(method string, url *Url) (*Request, error) {
	req, err := http.NewRequest("GET", url.RequestUrl(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", url.host)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "close")

	return &Request{Req: req, Url: url}, nil
}

func (req Request) RequestString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("connecting to %s\n", req.Url.host))
	sb.WriteString(fmt.Sprintf("Sending request GET /%s HTTP/1.1\n", req.Url.path))
	for k := range req.Req.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, req.Req.Header.Get(k)))
	}
	return sb.String()
}

func DumpResponse(resp *http.Response) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n", resp.Proto, resp.Status))
	sb.WriteString(fmt.Sprintf("Date: %s\n", resp.Header.Get("Date")))
	sb.WriteString(fmt.Sprintf("Content-Type: %s\n", resp.Header.Get("Content-Type")))
	sb.WriteString(fmt.Sprintf("Content-Length: %s\n", resp.Header.Get("Content-Length")))
	sb.WriteString(fmt.Sprintf("Connection: %s\n", resp.Header.Get("Connection")))
	sb.WriteString(fmt.Sprintf("Server: %s\n", resp.Header.Get("Server")))
	sb.WriteString(fmt.Sprintf("Access-Control-Allow-Origin: %s\n", resp.Header.Get("Access-Control-Allow-Origin")))
	sb.WriteString(fmt.Sprintf("Access-Control-Allow-Credentials: %s\n", resp.Header.Get("Access-Control-Allow-Credentials")))
	return sb.String()
}
