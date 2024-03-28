package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/ingest"
	v1 "github.com/hyperbolicresearch/hlog/web/api/v1"
)

func main() {
	log.Println("Hlog engine started...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// We load the configurations by reading the config.yaml, otherwise
	// (if it fails to load), we load the default configurations.
	cfg, err := config.FromYAML("config.yaml")
	if err != nil {
		cfg = &config.DefaultConfig
	}

	// Spinning up the components needed for livetailing and
	// the metrics/observables.
	livetail := v1.NewLiveTail(cfg.APIv1)
	observablestail := v1.NewObservablesTail(cfg)
	go livetail.Start(sigchan)
	go observablestail.Start(sigchan)

	// Start the API that is serves the observers
	apiServer := v1.New(cfg.APIv1)
	go apiServer.Start(sigchan)

	// We then start the ingesting processes in MongoDB and in ClickHouse.
	mongodbIngester := ingest.NewMongoDBIngester(cfg)
	go mongodbIngester.Start(sigchan)
	// clickhouseIngester := ingest.NewClickHouseIngester(cfg)
	// go clickhouseIngester.Start(sigchan)

	<-sigchan
}
