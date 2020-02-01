package res

// Common server event
type ServerEvent struct {
	View  string      `json:"view"`
	Error string      `json:"err"`
	Data  interface{} `json:"data"`
}

// General data from game session
type Data struct {
	Owner        PlayerInfo  `json:"owner"`
	SessionID    int         `json:"s_id"`
	PlayersCount int         `json:"p_count"`
	PlayerList   PlayersInfo `json:"p_list"`
}

type PlayersInfo struct {
	Players []PlayerInfo
}

type PlayerInfo struct {
	Index int    `json:"id"`
	Name  string `json:"name"`
}
