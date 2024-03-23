package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/core"
	"github.com/hyperbolicresearch/hlog/internal/kafkaservice"
	"github.com/hyperbolicresearch/hlog/pkg/logger"
	"golang.org/x/net/websocket"
)

// LiveTail is a real-time, bridge between Kafka and a logging medium
// that allows the observation as they are occuring of newly ingested
// messages.
func LiveTail(config *config.Livetail, sigchan chan os.Signal) {
	kw, err := kafkaservice.NewKafkaWorker(&config.KafkaConfigs)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()
	kw.SubscribeTopics(config.KafkaTopics)

	logger := logger.New(config.DefaultLevel, os.Stdout)

	// We build and spin up a new websocket server that is responsible
	// for basically adding and removing connections to the list of
	// writers of the logger.
	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		buf := make([]byte, 1024)
		for {
			_, err := ws.Read(buf)
			if err == io.EOF {
				if err := logger.RemoveWriter(ws); err != nil {
					panic(err)
				}
				continue
			}
			// We assume that all call to this endpoint constitues
			// a demand of connection.
			if err := logger.AddWriter(ws); err != nil {
				panic(err)
			}
			ws.Write([]byte("Connected"))
		}
	}))

	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%v", config.WebsocketPort), nil)
		if err != nil {
			panic(err)
		}
	}()

	ticker := time.NewTicker(time.Second)
	run := true
	for run {
		select {
		case <-ticker.C:
			logger.RLock()
			fmt.Println(len(logger.Writers))
			logger.RUnlock()
		case <-sigchan:
			log.Printf("Caught signal: %v", sigchan)
			run = false
		default:
			ev, err := kw.Consumer.ReadMessage(time.Duration(100) * time.Millisecond)
			if err != nil {
				continue
			}
			var l core.Log
			if err := json.Unmarshal(ev.Value, &l); err != nil {
				fmt.Printf("error unmarshalling value %v: %v", ev.Value, err)
			} else {
				err := logger.Log(l)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
