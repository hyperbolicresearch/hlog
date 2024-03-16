package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/utils"
)

func init() {
	file, _ := os.Open(".env")
	if file != nil {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// We load the configurations by reading the config.yaml, otherwise
	// (if it fails to load), we load the default configurations.
	cfg, err := config.FromYAML("config.yaml")
	if err != nil {
		cfg = &config.DefaultConfig
	}

	// TODO make a better welcome message here
	log.Println("Hlog live tail (experimental) up and running...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)


	utils.LiveTail(cfg.Livetail)
	<-sigchan
}
