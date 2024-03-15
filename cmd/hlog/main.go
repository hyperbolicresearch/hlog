package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/ingest"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

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

	// We then start the ingesting processes in MongoDB and in ClickHouse.
	mongodbIngester := ingest.NewMongoDBIngester(cfg)
	clickhouseIngester := ingest.NewClickHouseIngester(cfg)
	go mongodbIngester.Start(sigchan)
	go clickhouseIngester.Start(sigchan)

	<-sigchan
}
