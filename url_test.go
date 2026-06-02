package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrlSuccess(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *Url
	}{
		{
			"default port",
			"http://eu.httpbin.org/get",
			&Url{
				proto: HTTP,
				host:  "eu.httpbin.org",
				port:  80,
				path:  "get",
			},
		},
		{
			"without explicit path",
			"http://eu.httpbin.org",
			&Url{
				proto: HTTP,
				host:  "eu.httpbin.org",
				port:  80,
				path:  "/",
			},
		},
		{
			"explicit port",
			"http://eu.httpbin.org:42/get",
			&Url{
				proto: HTTP,
				host:  "eu.httpbin.org",
				port:  42,
				path:  "get",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(wrapName(tc.name), func(t *testing.T) {
			got, err := ParseUrl(tc.input)
			assert.Nil(t, err)
			assert.EqualValues(t, tc.want, got, "Input: %q", tc.input)
		})
	}
}

func TestParseUrlError(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"missing seperator", "asdf", ErrMissingSeperator},
		{"unsupported protocol", "asdf://", ErrUnsupportedProto},
		{"empty host", "http://", ErrEmptyHost},
		{"negative port", "http://www.example.com:-1521", ErrInvalidPort},
		{"too large port", "http://www.example.com:65536", ErrInvalidPort},
	}

	for _, tc := range testCases {
		t.Run(wrapName(tc.name), func(t *testing.T) {
			_, err := ParseUrl(tc.input)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestParsePort(t *testing.T) {
	type WantPort struct {
		port uint16
		rest string
	}

	testCases := []struct {
		name     string
		input    string
		wantPort WantPort
		wantErr  error
	}{
		{"port without path", "80", WantPort{port: 80}, nil},
		{"port with slash", "80/", WantPort{port: 80, rest: "/"}, nil},
		{"port with path", "80/path", WantPort{port: 80, rest: "/path"}, nil},
		{"max port", "65535", WantPort{port: 65535}, nil},
		{"negative port", "-5134", WantPort{}, ErrInvalidPort},
		{"port too large", "65536", WantPort{}, ErrInvalidPort},
	}

	for _, tc := range testCases {
		t.Run(wrapName(tc.name), func(t *testing.T) {
			port, rest, err := parsePort(tc.input)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.EqualValuesf(t, tc.wantPort,
				WantPort{port: port, rest: rest}, "Input: %s", tc.input)
		})
	}
}

func wrapName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}
