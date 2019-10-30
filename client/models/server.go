package models

type ServerListResponse []struct {
	Server Server `json:"server"`
}

type Subnet struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type Server struct {
	Cancelled    bool     `json:"cancelled"`
	Dc           string   `json:"dc"`
	Flatrate     bool     `json:"flatrate"`
	IP           []string `json:"ip"`
	PaidUntil    string   `json:"paid_until"`
	Product      string   `json:"product"`
	ServerIP     string   `json:"server_ip"`
	ServerName   string   `json:"server_name"`
	ServerNumber int      `json:"server_number"`
	Status       string   `json:"status"`
	Subnet       []Subnet `json:"subnet"`
	Throttled    bool     `json:"throttled"`
	Traffic      string   `json:"traffic"`
}
