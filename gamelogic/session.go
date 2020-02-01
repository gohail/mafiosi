package gamelogic

import (
	"github.com/gohail/mafiosi/metadata/req"
	"github.com/gohail/mafiosi/metadata/res"
	"github.com/gohail/mafiosi/metadata/view"
	"github.com/gohail/mafiosi/model"
	"go.uber.org/zap"
)

type GameSession struct {
	Owner        model.Player
	GameId       int
	PlayersCount int
	Players      []model.Player
	IsOpen       bool
	JoinListener chan interface{}
	GameStarter  chan interface{}
}

func NewGameSession(ow model.Player, gameId int, playersCount int, players []model.Player) *GameSession {
	return &GameSession{Owner: ow, GameId: gameId, PlayersCount: playersCount, Players: players, IsOpen: true,
		JoinListener: make(chan interface{}, 10),
		GameStarter:  make(chan interface{}),
	}
}

// Main Game method
func (s *GameSession) StartSession() {
	var stopper chan interface{}
	go waitAllPlayers(s, stopper)
	// Send game info at least for Session Owner
	s.JoinListener <- struct{}{}

	var action req.ActionReq

	for {
		if err := s.Owner.Conn.ReadJSON(&action); err != nil {
			zap.S().Error(err)
		} else {
			if "NEXT" == action.Action {
				zap.S().Infof("SESSION:%d adding new players is closed", s.GameId)
				close(stopper)
				break
			}
		}
	}

}

func waitAllPlayers(s *GameSession, stop <-chan interface{}) {
	for {
		select {
		case <-s.JoinListener:
			s.SendStartInfoToAll()
		case <-stop:
			zap.S().Infof("SESSION:%d stopped joiner goroutine", s.GameId)
			return
		}
	}

}

func (s *GameSession) SendStartInfoToAll() {
	evt := res.ServerEvent{
		View: view.OwnerStartInfo,
		Data: s.getSessionData(),
	}
	if err := s.Owner.Conn.WriteJSON(evt); err != nil {
		zap.S().Error(err)
	}

	evt = res.ServerEvent{
		View: view.PlayerStartInfo,
		Data: s.getSessionData(),
	}
	for _, p := range s.Players[1:] {
		if err := p.Conn.WriteJSON(evt); err != nil {
			zap.S().Error(err)
		}
	}
}
