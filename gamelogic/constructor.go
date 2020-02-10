package gamelogic

import (
	"github.com/gohail/mafiosi/metadata/res"
)

func (s *GameSession) getPlayersInfo() res.PlayersInfo {
	var list []res.PlayerInfo

	for _, p := range s.Players {
		list = append(list, p.ToPlayerInfo())
	}

	return res.PlayersInfo{Players: list}
}

func (s *GameSession) getSessionData() res.Data {
	return res.Data{
		Owner:        s.Owner.ToPlayerInfo(),
		SessionID:    s.GameId,
		PlayersCount: s.PlayersCount,
		PlayerList:   s.getPlayersInfo(),
	}
}
