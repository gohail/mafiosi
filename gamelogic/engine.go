package gamelogic

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"mafiosi/model"
	"mafiosi/utils"
)

var GameEngine = InitEngine()



type Engine struct {
	sessions map[int]*GameSession
}

func InitEngine() *Engine{
	return &Engine{sessions: make(map[int]*GameSession)}
}

// Create GameSession with uuid as int
func (e *Engine) CreateSession(owner model.Player, gameId int) {
	players := make([]model.Player, 0)
	players = append(players, owner)
	session := NewGameSession(owner, gameId, 1, players)
	zap.S().Infof("SESSION:%d created", session.GameId)
	e.sessions[gameId] = session
	session.StartSession()
}

func (e *Engine) JoinToSession(p model.Player, sessionId int) error {
	s, ok := e.sessions[sessionId]
	if !ok {
		return errors.New(fmt.Sprintf("game with ID:%d not found", sessionId))
	}

	if !s.IsOpen {
		return errors.New(fmt.Sprintf("game #%d already started", sessionId))
	}

	s.Players = append(s.Players, p)
	s.PlayersCount = len(s.Players)
	zap.S().Infof("SESSION:%d player %s joined", sessionId, p.Name)
	return nil
}

func (e *Engine) NewSessionUUID() int {
	ID := utils.GenerateID()
	_, ok := e.sessions[ID]
	if !ok {
		return ID
	}
	return e.NewSessionUUID()
}

// Check session id in pool
func (e *Engine) CheckSessionId(id int) bool {
	_, ok := e.sessions[id]
	return ok
}

// Check session id in pool
func (e *Engine) CheckUniqName(id int, name string) bool {
	s, ok := e.sessions[id]
	if ok {
		for _, p := range s.Players {
			if name == p.Name {
				return false
			}
		}
	}
	return true
}
