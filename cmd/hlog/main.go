package main

import (
	"github.com/hyperbolicresearch/hlog/internal/server"
)

func main() {
	s := server.NewServer(":8000")
	s.Start()
}
