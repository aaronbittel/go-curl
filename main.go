package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var verbose bool
var method string

var positionalArgs = []string{"url"}

func init() {
	const (
		defaultVerbose = false
		verboseUsage   = "Verbose output (dump headers)"
	)

	flag.BoolVar(&verbose, "verbose", defaultVerbose, verboseUsage)
	flag.BoolVar(&verbose, "v", defaultVerbose, verboseUsage+" (shorthand)")
	flag.StringVar(&method, "X", http.MethodGet, "Specify method")

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
