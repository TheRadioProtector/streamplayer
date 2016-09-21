package types

import "time"

type Teams struct {
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"_links"`
	Teams []struct {
		ID          int       `json:"_id"`
		Name        string    `json:"name"`
		Info        string    `json:"info"`
		DisplayName string    `json:"display_name"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Logo        string    `json:"logo"`
		Banner      string    `json:"banner"`
		Background  string    `json:"background"`
		Links       struct {
			Self string `json:"self"`
		} `json:"_links"`
	} `json:"teams"`
}
