package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
)

type Server struct {
	games    sync.Map
	waiting  *GameInstance
	Mux      *http.ServeMux
	WsServer *http.Server
	config   Config
	current  int
	total    int
}

type Config struct {
	Port int `default:"8080"`
}

type Response struct {
	GameID   string
	PlayerID int
}

func NewServer(c Config) *Server {
	return &Server{
		games:   sync.Map{},
		waiting: nil,
		config:  c,
	}
}

func (s *Server) Start() {
	sm := http.NewServeMux()
	endCh := make(chan string)

	sm.HandleFunc("/play", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		if s.waiting != nil {
			g := s.waiting
			s.SendResponse(writer, g.id, 1)
			s.games.Store(g.id, g)
			s.waiting = nil
			return
		}
		gi := NewGameInstance(endCh)
		s.current++
		s.total++
		s.waiting = gi
		s.games.Store(gi.id, gi)
		fmt.Printf("Creating new game %s. %d games in progress. %d games total\n", gi.id, s.current, s.total)
		s.SendResponse(writer, gi.id, 0)
	})

	sm.HandleFunc("/game", func(writer http.ResponseWriter, request *http.Request) {
		g, id, c := s.connect(writer, request)
		// pass the connection to the game itself
		g.ConnectPlayer(id, c)
	})

	sm.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _, c := s.connect(writer, request)
		// Simple ping loop that doesn't affect the game
		go func() {
			for {
				_, _, err := c.Read(context.Background())
				if err != nil {
					break
				}
				err = c.Write(context.Background(), websocket.MessageText, []byte(""))
				if err != nil {
					break
				}
			}
			c.Close(200, "error")
		}()
	})

	sm.HandleFunc("/watch", func(writer http.ResponseWriter, request *http.Request) {
		g, _, c := s.connect(writer, request)
		g.addSpectator(c)
	})

	// Start game management thread
	go func() {
		for {
			select {
			case id := <-endCh:
				s.games.Delete(id)
				s.current--
				if s.waiting != nil && s.waiting.id == id {
					s.waiting = nil
				}
				fmt.Printf("Game %s ended. %d games in progress, %d games played\n", id, s.current, s.total)
			}
		}
	}()
	fmt.Printf("Starting server on port %v\n", s.config.Port)

	server := http.Server{Handler: cors.Default().Handler(sm), Addr: fmt.Sprintf(":%d", s.config.Port)}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (s *Server) connect(writer http.ResponseWriter, request *http.Request) (*GameInstance, int, *websocket.Conn) {
	options := &websocket.AcceptOptions{OriginPatterns: []string{"*"}}
	ctx := context.Background()
	c, err := websocket.Accept(writer, request, options)
	if err != nil {
		panic(err)
	}
	cp := ConnectPayload{}
	//expect a client connecting
	_, b, err := c.Read(ctx)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &cp)
	if err != nil {
		panic(err)
	}
	// Find the game
	g, ok := s.games.Load(cp.GameID)
	if !ok {
		b, _ := json.Marshal(g)
		err = c.Write(ctx, websocket.MessageText, b)
		if err != nil {
			panic(err)
		}
	}
	return g.(*GameInstance), cp.PlayerID, c
}

func (s *Server) Stop() {
	s.WsServer.Close()
}

func (s *Server) SendResponse(w http.ResponseWriter, gameID string, playerID int) {
	r := Response{
		GameID:   gameID,
		PlayerID: playerID,
	}
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}
}
