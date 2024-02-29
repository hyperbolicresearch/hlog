package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyperbolicresearch/hlog/utils"

	"github.com/joho/godotenv"
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
	log.Println("Hlog live tail (experimental) up and running...")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	utils.LiveTail()
	<-sigchan
}
