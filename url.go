package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Url struct {
	proto Protocol
	host  string
	port  uint16
	path  string
}

var (
	ErrMissingSeperator = errors.New("missing \"://\" seperator")
	ErrUnsupportedProto = errors.New("unsupported protocol")
	ErrEmptyHost        = errors.New("empty host")
	ErrInvalidPort      = errors.New("invalid port")
)

// Parses a Url from a string
func ParseUrl(input string) (*Url, error) {
	proto, rest, err := parseProto(input)
	if err != nil {
		return nil, err
	}

	host, rest, err := parseHost(rest)
	if err != nil {
		return nil, err
	}

	port := proto.Port()
	if strings.HasPrefix(rest, ":") {
		port, rest, err = parsePort(rest[1:])
		if err != nil {
			return nil, err
		}
	}

	path := rest
	if len(rest) == 0 {
		path = "/"
	}
	if strings.HasPrefix(rest, "/") {
		path = path[1:]
	}

	return &Url{
		proto: proto,
		host:  host,
		port:  port,
		path:  path,
	}, nil
}

func (url Url) RequestUrl() string {
	portStr := ""
	if url.port != url.proto.Port() {
		portStr = fmt.Sprintf(":%d", url.port)
	}
	return fmt.Sprintf("%s://%s%s/%s", url.proto, url.host, portStr, url.path)
}

// Parses a leading protocol from input.
//
// It returns the parsed protocol, remaining input.
func parseProto(input string) (proto Protocol, rest string, err error) {
	protoName, rest, found := strings.Cut(input, "://")
	if !found {
		return "", "", ErrMissingSeperator
	}

	proto, err = ParseProtocol(protoName)
	if err != nil {
		return "", "", fmt.Errorf("%w: %q", ErrUnsupportedProto, protoName)
	}

	return
}

// Parses a leading host from input.
//
// It returns the parsed host and the remaining input.
// The returned rest includes the leading ':' or '/' if present.
func parseHost(input string) (host string, rest string, err error) {
	isStopChar := func(c byte) bool {
		return c == ':' || c == '/'
	}

	i := 0
	for ; i < len(input) && !isStopChar(input[i]); i++ {
	}

	if i == 0 {
		return "", "", ErrEmptyHost
	}

	return input[:i], input[i:], nil
}

// Parses a leading port from input.
//
// It returns the parsed port and the remaining input.
func parsePort(input string) (port uint16, rest string, err error) {
	isDigit := func(i byte) bool {
		return i >= '0' && i <= '9'
	}

	i := 0
	for ; i < len(input) && isDigit(input[i]); i++ {
	}

	portU64, err := strconv.ParseUint(input[:i], 10, 16)
	if err != nil {
		return 0, "", fmt.Errorf("%w: invalid port %q", ErrInvalidPort, input[:i])
	}

	return uint16(portU64), input[i:], nil
}
