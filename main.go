package main

import (
	"github.com/bitwurx/jrpc2"
)

func main() {
	InitDatabase()
	s := jrpc2.NewServer(":8888", "/rpc")
	// interface goes here
	NewApiV1(nil, s)
	s.Start()
}
