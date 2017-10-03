package main

import (
	"github.com/bitwurx/jrpc2"
)

// ApiV1 is the version 1 implementation of the rpc methods.
type ApiV1 struct {
	// model the priority queue database model.
	// queues A represetation of priority queues by key.
	model Model
}

// NewApiV1 returns a new api version 1 rpc api instance
func NewApiV1(model Model, s *jrpc2.Server) *ApiV1 {
	return nil
}
