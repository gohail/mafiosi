package req

const NameReq string = "NAME_REQ"

type NameForm struct {
	Name string `json:"name"`
}

type GameOption struct {
	PlayersSeq []int `json:"p_sequence"`
	MafNum     int   `json:"maf_number"`
	Cop        bool  `json:"cop"`
}

type IdForm struct {
	ID int `json:"game_id"`
}

type TemplateReq struct {
	Template string `json:"template"`
}

type ActionReq struct {
	Action string `json:"action"`
}
