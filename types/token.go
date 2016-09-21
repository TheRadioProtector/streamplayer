package types

type AccessToken struct {
	Token            string `json:"token"`
	Sig              string `json:"sig"`
	MobileRestricted bool   `json:"mobile_restricted"`
}
