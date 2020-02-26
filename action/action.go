package action

import (
	"errors"
	"fmt"
	"github.com/gohail/mafiosi/gamelogic"
	"github.com/gohail/mafiosi/metadata/req"
	"github.com/gohail/mafiosi/metadata/res"
	"github.com/gohail/mafiosi/metadata/view"
	"github.com/gohail/mafiosi/model"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

const (
	CreateGame = "CREATE_GAME"
	JoinGame   = "JOIN_GAME"
	Cancel     = "CANCEL"
	Next       = "NEXT"
	Other      = "OTHER"
)

var u = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnHandler(w http.ResponseWriter, r *http.Request) {
	zap.S().Infow("new connection:", "URL:", r.URL)
	c, err := u.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}
	startActionListener(c)
}

func startActionListener(c *websocket.Conn) {
	var errMsg string
	for {
		if err := c.WriteJSON(res.ServerEvent{
			View:  view.StartView,
			Error: errMsg,
			Data:  struct{}{},
		}); err != nil {
			zap.S().Errorf("[ActionListener] WriteJSON err: %v", err)
			break
		}

		if errMsg != "" {
			errMsg = ""
		}

		var actRq req.ActionReq

		if err := c.ReadJSON(&actRq); err != nil {
			zap.S().Errorf("[ActionListener] ReadJSON err: %v", err)
			break
		}

		switch actRq.Action {
		case CreateGame:
			zap.S().Info("USER CALL CREATE ACTION")
			createGame(c)
		case JoinGame:
			zap.S().Info("USER CALL JOIN ACTION")
			joinGame(c)
			return
		case Cancel:
			zap.S().Info("USER CALL ABORT")
			return
		default:
			zap.S().Infof("UNKNOWN ACTION: \"%s\"", actRq.Action)
			errMsg = fmt.Sprintf("Invalid action: %s", actRq.Action)
		}
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

func joinGame(c *websocket.Conn) {
	gameId, err := gameIdReqLoop(c)
	if err != nil {
		zap.S().Errorf("join game err: ", err)
		return
	}
	name, err := playerNameReqLoop(c, gameId)
	if err != nil {
		zap.S().Errorf("join game err: ", err)
		return
	}
	p := model.Player{
		Conn:     c,
		PlayerId: 0,
		Name:     name,
		Role:     "",
		IsAlive:  true,
	}

	if err = gamelogic.GameEngine.JoinToSession(p, gameId); err != nil {
		zap.S().Error(err)
	}
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

		var body struct {
			req.IdForm
			req.ActionReq
		}

		if err := c.ReadJSON(&body); err != nil {
			zap.S().Error(err)
		}

		if body.Action == Cancel {
			return 0, errors.New("CANCEL ACTION SELECTED")
		}

		if gamelogic.GameEngine.CheckSessionId(body.ID) {
			return body.ID, nil
		} else {
			errMsg = fmt.Sprintf("Invalid game ID: %d", body.ID)
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

		var body struct {
			req.ActionReq
			req.NameForm
		}
		if err := c.ReadJSON(&body); err != nil {
			zap.S().Error(err)
			return "", err
		}

		if body.Action == Cancel {
			return "", errors.New("CANCEL ACTION SELECTED")
		}

		if body.Name != "" {
			if gamelogic.GameEngine.CheckUniqName(sId, body.Name) {
				return body.Name, nil
			} else {
				errMsg = "select a unique name"
			}
		} else {
			errMsg = "name can't be empty"
			zap.S().Info("user typo empty name!")
		}
	}
}
