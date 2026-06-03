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

func DumpResponse(resp *http.Response) {
	fmt.Println(resp.Proto, resp.Status)
	fmt.Println("Date:", resp.Header.Get("Date"))
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))
	fmt.Println("Content-Length:", resp.Header.Get("Content-Length"))
	fmt.Println("Connection:", resp.Header.Get("Connection"))
	fmt.Println("Server:", resp.Header.Get("Server"))
	fmt.Println("Access-Control-Allow-Origin:", resp.Header.Get("Access-Control-Allow-Origin"))
	fmt.Println("Access-Control-Allow-Credentials:", resp.Header.Get("Access-Control-Allow-Credentials"))
}
