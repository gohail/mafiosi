package req

const NameReq string = "NAME_REQ"

type NameForm struct {
	Name string `json:"name"`
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
