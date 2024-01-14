package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	Name      string
	Id        uuid.UUID
	IsRunning bool
	QuitC     chan struct{}
	sync.RWMutex
	*http.Server
}

// NewServer creates and returns a new Server instance.
func NewServer(addr string) *Server {

	r := CreateAndConfigureMultiplexer()
	http_s := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	s := Server{
		Name:   "server",
		Id:     uuid.New(),
		Server: http_s,
	}

	return &s
}

// CreateAndConfigureMultiplexer creates a new mux.Router struct, add
// the needed handling functions according to the needs of the API.
func CreateAndConfigureMultiplexer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/log", LogHandler)
	return r
}

// Start starts the server and blocks until it it receives a quit signal
// from the QuitC channel.
func (s *Server) Start() {
	s.Lock()
	defer s.Unlock()
	s.IsRunning = true

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println("An error occured while launching the server.")
		}
		s.Lock()
		s.IsRunning = false
		s.Unlock()
	}()

	log.Println("Server started...")
	<- s.QuitC
}
