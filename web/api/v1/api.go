package v1

import (
	"io"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"

	"github.com/hyperbolicresearch/hlog/config"
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

// HandleLiveInit read the last
func (s *Server) HandleLiveInit(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) HandleGenObservables(w http.ResponseWriter, r *http.Request) {}
func (s *Server) HandleObserve(w http.ResponseWriter, r *http.Request)        {}
func (s *Server) HandleInfo(w http.ResponseWriter, r *http.Request)           {}
