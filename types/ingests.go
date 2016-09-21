package types

type Ingests struct {
	Links struct {
		Self string `json:"self"`
	} `json:"_links"`
	Ingests []struct {
		Name         string  `json:"name"`
		Availability float64 `json:"availability"`
		ID           int     `json:"_id"`
		Default      bool    `json:"default"`
		URLTemplate  string  `json:"url_template"`
	} `json:"ingests"`
}
