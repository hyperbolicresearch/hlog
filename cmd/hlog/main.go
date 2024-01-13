package main

import (
	"fmt"

	"github.com/hyperbolicresearch/hlog/internal/core"
)

func main() {
	msg := "hello, there!"
	data := struct{ foo string }{foo: "bar"}
	w := "https://hyperbolicresearch.com"
	log := core.NewLog(msg, data, w, core.LogOptions{})
	fmt.Println(log.String())
}
