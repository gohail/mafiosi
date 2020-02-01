package model

import (
	"github.com/gohail/mafiosi/metadata/res"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	PlayerId int
	Name     string
	Role     string
	IsAlive  bool
}

func NewPlayer(c *websocket.Conn, id int, name string, role string) *Player {
	return &Player{
		Conn:     c,
		PlayerId: id,
		Name:     name,
		Role:     role,
		IsAlive:  true,
	}
}

func (p *Player) ToPlayerInfo() res.PlayerInfo {
	return res.PlayerInfo{
		Index: p.PlayerId,
		Name:  p.Name,
	}
}
