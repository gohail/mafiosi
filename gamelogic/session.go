package gamelogic

import (
	"errors"
	"fmt"
	"github.com/gohail/mafiosi/action"
	"github.com/gohail/mafiosi/metadata/req"
	"github.com/gohail/mafiosi/metadata/res"
	"github.com/gohail/mafiosi/metadata/view"
	"github.com/gohail/mafiosi/model"
	"github.com/gorilla/websocket"
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
func (s *GameSession) PrepareSession() {
	var stopper chan interface{}
	go waitJoinPlayers(s, stopper)
	// Send game info at least for Session Owner
	s.JoinListener <- struct{}{}

	// Waiting from owner NEXT actionReq
	if err := waitAction(action.Next, s.Owner.Conn); err != nil {
		zap.S().Error("FATAL SESSION ERROR: exit from game session!")
		zap.S().Error(err)
		return
	}

	//close session for new players
	s.IsOpen = false
	zap.S().Infof("SESSION:%d access to the game is closed, waiting game-options from owner.")

	//Send to owner game-option view
	s.sendToPlayer(s.Owner, res.ServerEvent{
		View: "GAME_OPTION",
		Data: s.getSessionData(),
	})

	//Send to other players wait view
	s.sendToJoinPlayers(res.ServerEvent{
		View: "WAITING_FOR_START",
		Data: s.getSessionData(),
	})

	var opt req.GameOption
	// Waiting for owner GameOption action
	if err := s.Owner.Conn.ReadJSON(&opt); err != nil {
		zap.S().Error("FATAL SESSION ERROR: exit from game session!")
		zap.S().Error(err)
		return
	}

	s.clearRes()
}

func waitJoinPlayers(s *GameSession, stop <-chan interface{}) {
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
	data := s.getSessionData()
	evt := res.ServerEvent{
		View: view.OwnerStartInfo,
		Data: data,
	}
	if err := s.Owner.Conn.WriteJSON(evt); err != nil {
		zap.S().Error(err)
	}

	evt = res.ServerEvent{
		View: view.PlayerStartInfo,
		Data: data,
	}
	for _, p := range s.Players[1:] {
		if err := p.Conn.WriteJSON(evt); err != nil {
			zap.S().Error(err)
		}
	}
}

func (s *GameSession) sendToAll(mess interface{}) {
	for _, p := range s.Players {
		if err := p.Conn.WriteJSON(mess); err != nil {
			zap.S().Error(err)
		}
	}
}

func (s *GameSession) sendToPlayer(p model.Player, mess interface{}) {
	if err := p.Conn.WriteJSON(mess); err != nil {
		zap.S().Error(err)
	}
}

func (s *GameSession) sendToJoinPlayers(mess interface{}) {
	for _, p := range s.Players[1:] {
		if err := p.Conn.WriteJSON(mess); err != nil {
			zap.S().Error(err)
		}
	}
}

func waitAction(action string, c *websocket.Conn) error {
	var mess req.ActionReq
	// Waiting for owner NEXT action
	if err := c.ReadJSON(&mess); err != nil {
		return err
	}
	if action == mess.Action {
		//zap.S().Infof("SESSION:%d adding new players is closed", s.GameId)
		//close(stopper)
		return nil
	}

	return errors.New(fmt.Sprintf(`Expect action: "%s" actual: "%s"`, action, mess.Action))
}

func (s *GameSession) clearRes() {
	//TODO
}
