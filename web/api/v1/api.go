package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/net/websocket"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/core"
	"github.com/hyperbolicresearch/hlog/internal/kafkaservice"
	"github.com/hyperbolicresearch/hlog/pkg/logger"
)

type Server struct {
	sync.RWMutex
	*http.Server
	Config               *config.Config
	WebsocketConnections []*websocket.Conn
	Logger               *logger.Logger
}

func New(config *config.Config) *Server {
	srv := &Server{
		Config:               config,
		WebsocketConnections: make([]*websocket.Conn, 0, config.Livetail.MaxWebsocketConnections),
	}
	return srv
}

// Configure will create the server and register the multiplexer.
func (s *Server) Configure() error {
	mux := http.NewServeMux()

	mux.Handle("/live", websocket.Handler(s.HandleLive))
	mux.Handle("/liveinit", http.HandlerFunc(s.HandleLiveInit))
	mux.Handle("/genericobservables", http.HandlerFunc(s.HandleGenObservables))
	mux.Handle("/observe/{observable_id}", http.HandlerFunc(s.HandleObserve))
	mux.Handle("/info", http.HandlerFunc(s.HandleInfo))

	s.Server = &http.Server{
		Addr:    s.Config.API.ServerAddr,
		Handler: mux,
	}
	return nil
}

// Start will start the http server and start listening fro incoming requests.
func (s *Server) Start(sigchan chan os.Signal) {
	s.Configure()
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	<-sigchan
}

// Stop should stop all opened websocket connections. It is intended to be called
// when Server.Shutdown() is called since the latter will not take care of closing
// hijacked connections such as websocket's.
func (s *Server) Stop() error {

	return nil
}

// HandleLive connects to Kafka and pipes the incomming messages to
// the connected "clients"
func (s *Server) HandleLive(ws *websocket.Conn) {
	// If the websocket connection is already in, move on
	s.RLock()
	for _, v := range s.WebsocketConnections {
		if v == ws {
			s.RUnlock()
			return
		}
	}
	s.RUnlock()

	s.Lock()
	s.WebsocketConnections = append(s.WebsocketConnections, ws)
	s.Unlock()

	buf := make([]byte, 1024)
	for {
		_, err := ws.Read(buf)
		if err := s.Config.Livetail.Logger.AddWriter(ws); err != nil {
			panic(err)
		}
		if err == io.EOF {
			s.Lock()
			for i, v := range s.WebsocketConnections {
				if v == ws {
					s.WebsocketConnections = append(s.WebsocketConnections[:i], s.WebsocketConnections[i+1:]...)
				}
			}
			s.Unlock()
			if err := s.Config.Livetail.Logger.RemoveWriter(ws); err != nil {
				panic(err)
			}
			break
		}
	}
}

// HandleLiveInit read the last config.Livetail.InitLogsLoadedCount logs from
// Kafka and returns them to the caller.
func (s *Server) HandleLiveInit(w http.ResponseWriter, r *http.Request) {
	// We make dedicated Kafka configurations because we need a brand new
	// way to interact with the Kafka server. For instance, we don't want to
	// commit messages, since we want to read and reread each time, we also
	// want to start reading from the tail (latest)
	kafkaConfig := config.Kafka{
		Server:           s.Config.Livetail.KafkaConfigs.Server,
		GroupId:          "hlog-livetail-init-default",
		AutoOffsetReset:  "latest",
		EnableAutoCommit: false,
	}
	kw, err := kafkaservice.NewKafkaWorker(&kafkaConfig)
	if err != nil {
		panic(err)
	}
	kw.ConfigureConsumer()

	err = kw.Consumer.Assign(
		[]kafka.TopicPartition{{
			Topic:     &s.Config.Livetail.KafkaTopics[0],
			Partition: 0,
			// TODO: Fix this, it should not be hardcoded
			Offset: kafka.OffsetTail(100),
		}})
	if err != nil {
		panic(err)
	}

	// TODO benchmark this to evaluate how the marshalling and
	// the unmarshallign penalize the process in term of performance.
	values := make([]*core.Log, 0, s.Config.Livetail.InitLogsLoadedCount)
	for i := 0; i < s.Config.Livetail.InitLogsLoadedCount; i++ {
		msg, err := kw.Consumer.ReadMessage(-1)
		if err != nil {
			panic(err)
		}
		var l core.Log
		if err := json.Unmarshal(msg.Value, &l); err != nil {
			panic(err)
		}
		values = append(values, &l)
	}

	slices.Reverse(values)
	jsonValues, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonValues))
}

// HandleGenObservables returns the general observables for the dashboard.
func (s *Server) HandleGenObservables(w http.ResponseWriter, r *http.Request) {

}

// TODO -----------------------------------------------------------------
func (s *Server) HandleObserve(w http.ResponseWriter, r *http.Request) {}
func (s *Server) HandleInfo(w http.ResponseWriter, r *http.Request)    {}
