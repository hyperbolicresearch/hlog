package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperbolicresearch/hlog/config"
	utils "github.com/hyperbolicresearch/hlog/utils"
)

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// We load the configurations by reading the config.yaml, otherwise
	// (if it fails to load), we load the default configurations.
	cfg, err := config.FromYAML("config.yaml")
	if err != nil {
		cfg = &config.DefaultConfig
	}

	fmt.Println("Producer simulator started...")
	utils.GenerateRandomLogs(cfg, sigchan)
}
