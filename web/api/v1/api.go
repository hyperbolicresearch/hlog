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
)

// Server is the server that runs and exposes the API
type Server struct {
	sync.RWMutex
	*http.Server
	config                       *config.APIv1
	liveTailWebsocketConnections []*websocket.Conn
	genObsWebsocketConnections   []*websocket.Conn // TODO: Maybe we can just use one ????
}

// New creates and returns a new API server instance.
func New(config *config.APIv1) *Server {
	srv := &Server{
		config:                       config,
		liveTailWebsocketConnections: make([]*websocket.Conn, 0, config.MaxLiveTailWebsocketConnections),
		genObsWebsocketConnections:   make([]*websocket.Conn, 0, config.MaxGenObsWebsocketConnections),
	}
	return srv
}

// Configure will create the HTTP server component of the API server
// and register its multiplexer.
func (s *Server) Configure() error {
	mux := http.NewServeMux()

	mux.Handle("/live", websocket.Handler(s.HandleLive))
	mux.Handle("/liveinit", http.HandlerFunc(s.HandleLiveInit))
	mux.Handle("/genericobservables", websocket.Handler(s.HandleGenObservables))
	mux.Handle("/observe/{observable_id}", http.HandlerFunc(s.HandleObserve))
	mux.Handle("/info", http.HandlerFunc(s.HandleInfo))

	s.Server = &http.Server{
		Addr:    s.config.ServerAddr,
		Handler: mux,
	}
	return nil
}

// Start will start the API server.
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

// Stop should stop all opened websocket connections. It is intended to 
// be called when Server.Shutdown() is called since the latter will not
// take care of closing hijacked connections such as websocket's.
func (s *Server) Stop() error {
	// TODO implement
	return nil
}

// HandleLive connects to Kafka and pipes the incomming messages to
// the connected "clients"
func (s *Server) HandleLive(ws *websocket.Conn) {
	// If the websocket connection is already in, move on
	s.RLock()
	for _, v := range s.liveTailWebsocketConnections {
		if v == ws {
			s.RUnlock()
			return
		}
	}
	s.RUnlock()

	s.Lock()
	s.liveTailWebsocketConnections = append(s.liveTailWebsocketConnections, ws)
	s.Unlock()

	buf := make([]byte, 1024)
	for {
		_, err := ws.Read(buf)
		if err == io.EOF {
			s.Lock()
			for i, v := range s.liveTailWebsocketConnections {
				if v == ws {
					s.liveTailWebsocketConnections = append(
						s.liveTailWebsocketConnections[:i],
						s.liveTailWebsocketConnections[i+1:]...)
				}
			}
			s.Unlock()
			if err := s.config.LivetailLogger.RemoveWriter(ws); err != nil {
				panic(err)
			}
			break
		}
		if err := s.config.LivetailLogger.AddWriter(ws); err != nil {
			panic(err)
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
		Server:           s.config.KafkaConfigs.Server,
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
			Topic:     &s.config.KafkaTopics[0],
			Partition: 0,
			// TODO: Fix this, it should not be hardcoded
			Offset: kafka.OffsetTail(100),
		}})
	if err != nil {
		panic(err)
	}

	// TODO benchmark this to evaluate how the marshalling and
	// the unmarshallign penalize the process in term of performance.
	values := make([]*core.Log, 0, s.config.InitLogsLoadedCount)
	for i := 0; i < s.config.InitLogsLoadedCount; i++ {
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
func (s *Server) HandleGenObservables(ws *websocket.Conn) {
	// If the conn is already there, move on
	s.RLock()
	for _, v := range s.genObsWebsocketConnections {
		if v == ws {
			s.RUnlock()
			return
		}
	}
	s.RUnlock()

	s.Lock()
	s.genObsWebsocketConnections = append(s.genObsWebsocketConnections, ws)
	s.Unlock()

	buf := make([]byte, 1024)
	for {
		if _, err := ws.Read(buf); err != nil {
			if err == io.EOF {
				s.Lock()
				for i, v := range s.genObsWebsocketConnections {
					if v == ws {
						s.genObsWebsocketConnections = append(
							s.genObsWebsocketConnections[:i],
							s.genObsWebsocketConnections[i+1:]...)
					}
				}
				s.Unlock()
				err := s.config.GeneralObservablesLogger.RemoveWriter(ws)
				if err != nil {
					panic(err)
				}
				break
			}
		}
		if err := s.config.GeneralObservablesLogger.AddWriter(ws); err != nil {
			panic(err)
		}
	}
}

// TODO -----------------------------------------------------------------
func (s *Server) HandleObserve(w http.ResponseWriter, r *http.Request) {}
func (s *Server) HandleInfo(w http.ResponseWriter, r *http.Request)    {}
