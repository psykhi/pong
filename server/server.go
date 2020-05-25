package server

import (
	"context"
	"encoding/json"
	"github.com/rs/cors"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
)
import "github.com/satori/go.uuid"

type Server struct {
	games    sync.Map
	waiting  map[string]*GameInstance
	Mux      *http.ServeMux
	WsServer *http.Server
}

type Response struct {
	GameID   string
	PlayerID int
}

const PORT = ":8080"
const ADDRESS = "pong-wasm.ew.r.appspot.com"
const HTTP_SCHEME = "https://"
const WS_SCHEME = "wss"

func NewServer() *Server {
	return &Server{
		games:   sync.Map{},
		waiting: map[string]*GameInstance{},
	}
}

func (s *Server) Start() {
	sm := http.NewServeMux()

	sm.HandleFunc("/play", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("content-type", "application/json")
		if len(s.waiting) != 0 {
			for id, g := range s.waiting {
				s.SendResponse(writer, id, 1)
				delete(s.waiting, id)
				s.games.Store(id, g)
				return
			}
		}
		gi := NewGameInstance()
		id := uuid.NewV4()
		s.waiting[id.String()] = gi
		s.games.Store(id.String(), gi)
		s.SendResponse(writer, id.String(), 0)
	})
	options := &websocket.AcceptOptions{OriginPatterns: []string{"*"}}
	sm.HandleFunc("/game", func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		c, err := websocket.Accept(writer, request, options)
		if err != nil {
			panic(err)
		}
		cp := ConnectPayload{}
		//expect a render connecting
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
			return

		}
		// pass the connection to the game itself
		g.(*GameInstance).ConnectPlayer(cp.PlayerID, c)
	})

	sm.HandleFunc("/watch", func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		c, err := websocket.Accept(writer, request, options)
		if err != nil {
			panic(err)
		}
		cp := ConnectPayload{}
		//expect a render connecting
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
		g.(*GameInstance).addSpectator(c)
	})

	server := http.Server{Handler: cors.Default().Handler(sm), Addr: PORT}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	// TODO handle end of game LOL
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
