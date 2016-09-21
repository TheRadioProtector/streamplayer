package types

type Root struct {
	Identified bool `json:"identified"`
	Links      struct {
		User    string `json:"user"`
		Channel string `json:"channel"`
		Search  string `json:"search"`
		Streams string `json:"streams"`
		Ingests string `json:"ingests"`
		Teams   string `json:"teams"`
	} `json:"_links"`
	Token struct {
		Valid         bool        `json:"valid"`
		Authorization interface{} `json:"authorization"`
	} `json:"token"`
}
