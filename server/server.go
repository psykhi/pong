package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/psykhi/pong/game"
	"net/http"
)
import "github.com/satori/go.uuid"

type Server struct {
	games    map[string]*GameInstance ``
	waiting  map[string]*GameInstance
	Mux      *http.ServeMux
	WsServer *http.Server
}

type Response struct {
	GameID   string
	PlayerID int
}

func NewServer() *Server {
	return &Server{
		games:   map[string]*GameInstance{},
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
				s.games[id] = g
				return
			}
		}
		g := game.NewState()
		id := uuid.NewV4()
		s.waiting[id.String()] = &GameInstance{State: g}
		s.SendResponse(writer, id.String(), 0)
	})
	sm.HandleFunc("/game", func(writer http.ResponseWriter, request *http.Request) {
		up := &websocket.Upgrader{}
		c, err := up.Upgrade(writer, request, nil)
		if err != nil {
			panic(err)
		}
		cp := ConnectPayload{}
		//expect a client connecting
		err = c.ReadJSON(&cp)
		if err != nil {
			panic(err)
		}
		// Find the game
		g, ok := s.games[cp.GameID]
		if !ok {
			err = c.WriteJSON(ok)
			if err != nil {
				panic(err)
			}
		}
		// pass the connection to the game itself
		g.ConnectPlayer(cp.PlayerID, c)
	})

	sm.HandleFunc("/watch", func(writer http.ResponseWriter, request *http.Request) {
		up := &websocket.Upgrader{}
		c, err := up.Upgrade(writer, request, nil)
		if err != nil {
			panic(err)
		}
		cp := ConnectPayload{}
		//expect a client connecting
		err = c.ReadJSON(&cp)
		if err != nil {
			panic(err)
		}
		// Find the game
		g, ok := s.games[cp.GameID]
		if !ok {
			err = c.WriteJSON(ok)
			if err != nil {
				panic(err)
			}
		}
		g.addSpectator(c)
	})

	server := http.Server{Handler: sm, Addr: ":3010"}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
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
