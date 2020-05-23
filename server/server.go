package server

import (
	"encoding/json"
	"github.com/psykhi/pong/game"
	"net/http"
)
import "github.com/satori/go.uuid"

type Server struct {
	games   map[string]*game.State
	waiting map[string]*game.State
}

type Response struct {
	GameID   string
	PlayerID int
}

func NewServer() *Server {
	return &Server{
		games:   map[string]*game.State{},
		waiting: map[string]*game.State{},
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
		s.waiting[id.String()] = g
		s.waiting[id.String()] = g
		s.SendResponse(writer, id.String(), 0)
	})

	err := http.ListenAndServe(":3010", sm)
	if err != nil {
		panic(err)
	}
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
