package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mafiosi/dataform/req"
	"mafiosi/dataform/res"
	"mafiosi/dataform/view"
	"mafiosi/gamelogic"
	"mafiosi/model"
	"net/http"
)

func stream(w http.ResponseWriter, r *http.Request) {
	zap.S().Infow("new connection:", "URL:", r.URL)
	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Error(err)
	}

	if err = c.WriteJSON(res.ServerEvent{
		View: view.ConnectionAction,
	}); err != nil {
		zap.S().Error(err)
	}

	var action req.ActionReq
	err = c.ReadJSON(&action)
	if err != nil {
		zap.S().Error(err)
	}

	switch action.Action {
	case "CREATE":
		zap.S().Info("USER CALL CREATE ACTION")
		createGame(c)
	case "JOIN":
		zap.S().Info("USER CALL JOIN ACTION")
		joinAction(c)
	default:
		zap.S().Infof("Unknown ACTION: %s", action.Action)
		if err = c.Close(); err != nil {
			zap.Error(err)
		}
	}
}

func joinAction(c *websocket.Conn) {
	gameId, err := gameIdReqLoop(c)
	if err != nil {
		zap.S().Error(err)
		return
	}
	name, err := playerNameReqLoop(c, gameId)
	if err != nil {
		zap.S().Error(err)
		return
	}
		p := model.Player{
		Conn:    c,
		PlayerId:  0,
		Name:    name,
		Role:    "",
		IsAlive: true,
	}

	if err = gamelogic.GameEngine.JoinToSession(p, gameId); err != nil{
		zap.S().Error(err)
	}
}

func createGame(c *websocket.Conn) {
	sID := gamelogic.GameEngine.NewSessionUUID()

	name, err := playerNameReqLoop(c, 0)
	if err != nil {
		zap.S().Errorf("create game err: ", err)
		return
	}

	zap.S().Infof("player %s create new game", name)

	u := model.Player{
		Conn:     c,
		PlayerId: 0,
		Name:     name,
		Role:     "",
		IsAlive:  true,
	}

	gamelogic.GameEngine.CreateSession(u, sID)
}

func gameIdReqLoop(c *websocket.Conn) (id int, err error) {
	var errMsg string
	for {
		if err := c.WriteJSON(res.ServerEvent{
			View:  view.ReqGameId,
			Error: errMsg,
		}); err != nil {
			return 0, err
		}

		var gameId req.IdForm
		if err := c.ReadJSON(&gameId); err != nil {
			zap.S().Error(err)
		}

		if gamelogic.GameEngine.CheckSessionId(gameId.ID) {
			return gameId.ID, nil
		} else {
			errMsg = fmt.Sprintf("Invalid game ID: %d", gameId.ID)
			zap.S().Error(errMsg)
		}
	}
}

func playerNameReqLoop(c *websocket.Conn, sId int) (string, error) {
	var errMsg string
	for {
		if err := c.WriteJSON(res.ServerEvent{
			View:  view.ReqName,
			Error: errMsg,
		}); err != nil {
			return "", err
		}

		var name req.NameForm
		if err := c.ReadJSON(&name); err != nil {
			zap.S().Error(err)
			return "", err
		}

		if name.Name != "" {
			if gamelogic.GameEngine.CheckUniqName(sId, name.Name) {
				return name.Name, nil
			} else {
				errMsg = "select a unique name"
			}
		} else {
			errMsg = "name can't be empty"
			zap.S().Info("user typo empty name!")
		}
	}
}
