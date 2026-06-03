package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s <url>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
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

	req, err := NewRequest("GET", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: creating request: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(req.RequestString())

	client := &http.Client{}
	resp, err := client.Do(req.Req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: sending request %s: %s\n", url.RequestUrl(), err)
		os.Exit(1)
	}
	DumpResponse(resp)
	fmt.Println()

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: reading response body %s: %s\n", url.RequestUrl(), err)
		os.Exit(1)
	}

	fmt.Println(string(body))
}
