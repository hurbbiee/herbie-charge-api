package dto

type ChargeRequest struct {
	IDToken string `json:"idToken"`
	Amount  int    `json:"amount"`
}

type ChargeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
