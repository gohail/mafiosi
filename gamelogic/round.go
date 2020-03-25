package gamelogic

import "github.com/gohail/mafiosi/model"

type Round interface {
	Start([]model.Player) error
}
