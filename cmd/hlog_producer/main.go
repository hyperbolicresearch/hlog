package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	utils "github.com/hyperbolicresearch/hlog/utils"
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
	fmt.Println("Producer simulator started...")
	stop := make(chan struct{})
	utils.GenerateRandomLogs(stop)
	<-stop
}
