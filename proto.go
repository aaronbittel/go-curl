package main

import (
	"fmt"
	"strings"
)

type Protocol string

const (
	HTTP Protocol = "http"
)

func (p Protocol) String() string {
	return string(p)
}

func (p Protocol) Port() uint16 {
	switch p {
	case HTTP:
		return 80
	default:
		return 0
	}
}

func ParseProtocol(proto string) (Protocol, error) {
	switch strings.ToLower(proto) {
	case "http":
		return HTTP, nil
	default:
		return "", fmt.Errorf("Unknown protocol: %s", proto)
	}
}
