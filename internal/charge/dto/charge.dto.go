package dto

type ChargeRequest struct {
	IDToken string `json:"idToken"`
	Amount  int    `json:"amount"`
}

type ChargeResult struct {
	TransactionID string `json:"transactionId"`
	Amount        int    `json:"amount"`
	Balance       int    `json:"balance"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
}

type ChargeResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *ChargeResult `json:"data,omitempty"`
}
