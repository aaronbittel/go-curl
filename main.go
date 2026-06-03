package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type headers http.Header

func (h headers) String() string {
	var sb strings.Builder
	i := 0
	for k, v := range h {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", k, v))
		i++
	}
	return sb.String()
}

func (h headers) Set(arg string) error {
	key, value, found := strings.Cut(arg, ":")
	if !found {
		return errors.New("missing \":\" in header")
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	http.Header(h).Add(key, value)
	return nil
}

var verbose bool
var method string
var data string
var headerFlag = make(headers)

var positionalArgs = []string{"url"}

func init() {
	const (
		defaultVerbose = false
		verboseUsage   = "Verbose output (dump headers)"

		headerUsage = "Add a header (can be used mulite times)"
	)

	flag.BoolVar(&verbose, "verbose", defaultVerbose, verboseUsage)
	flag.BoolVar(&verbose, "v", defaultVerbose, verboseUsage+" (shorthand)")
	flag.StringVar(&method, "X", http.MethodGet, "Specify method")
	flag.Var(&headerFlag, "header", headerUsage)
	flag.Var(&headerFlag, "H", headerUsage+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s <url>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() > len(positionalArgs) {
		fmt.Fprintf(os.Stderr,
			"WARNING: all flags must be provided before the positional args: %q",
			strings.Join(positionalArgs, ", "))
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	urlInput := flag.Arg(0)
	url, err := ParseUrl(urlInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	req := NewRequest(method, url)
	for k, vs := range headerFlag {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	if verbose {
		printOutgoing(req.RequestString())
	}

	resp, err := req.Send()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sending request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if verbose {
		printIncoing(DumpResponse(resp))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: reading body: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(bytes.TrimSpace(data)))
}

func printOutgoing(str string) {
	printVerbose(str, true)
}

func printIncoing(str string) {
	printVerbose(str, false)
}

func printVerbose(str string, outgoing bool) {
	prefix := "< "
	if outgoing {
		prefix = "> "
	}

	for _, line := range strings.Split(str, "\n") {
		fmt.Println(prefix + line)
	}
}
