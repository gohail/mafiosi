package model

import (
	"github.com/gohail/mafiosi/metadata/res"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	PlayerId int
	Name     string
	Role     Role
	IsAlive  bool
}

func (p *Player) ToPlayerInfo() res.PlayerInfo {
	return res.PlayerInfo{
		Index: p.PlayerId,
		Name:  p.Name,
	}
}
