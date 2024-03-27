package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/utils"
	v1 "github.com/hyperbolicresearch/hlog/web/api/v1"
)

func main() {
	// We load the configurations by reading the config.yaml, otherwise
	// (if it fails to load), we load the default configurations.
	// We load the configurations by reading the config.yaml, otherwise
	// (if it fails to load), we load the default configurations.
	cfg, err := config.FromYAML("config.yaml")
	if err != nil {
		cfg = &config.DefaultConfig
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Launching the livetail. But why ???
	// utils.Livetail is basically a Kafka consumer which reads messages
	// as they arrive to the system.
	go utils.LiveTail(cfg.Livetail, sigchan)

	log.Println("hlogAPI up and running...")
	server := v1.New(cfg)
	server.Start(sigchan)
}
