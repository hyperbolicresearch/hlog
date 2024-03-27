package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperbolicresearch/hlog/config"
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

	// Spinning up the components needed for livetailing and
	// the metrics/observables.
	livetail := v1.NewLiveTail(cfg.APIv1)
	observablestail := v1.NewObservablesTail(cfg)
	go livetail.Start(sigchan)
	go observablestail.Start(sigchan)

	log.Println("hlogAPI up and running...")
	server := v1.New(cfg.APIv1)
	server.Start(sigchan)
}
