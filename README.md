# Build Your Own curl

Implement a basic version of the Unix curl utility in Go, following the [Build Your Own
curl](https://codingchallenges.fyi/challenges/challenge-curl).

## Usage

```bash
go build -o build/go-curl
./build/go-curl http://eu.httpbin.org/get
```

- using Taskfile

```bash
task run -- http://eu.httpbin.org/get
```
