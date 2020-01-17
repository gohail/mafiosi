package gamelogic

import (
	"fmt"
	"mafiosi/dataform/res"
)

func (s *GameSession) getPlayersInfo() res.PlayersInfo {
	var list []res.PlayerInfo

	for _, p := range s.Players {
		fmt.Println("add NEW PLAYER INFO")
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
