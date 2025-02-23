package models

// TransactionRequest a standard request structure for the transactions
type TransactionRequest struct {
	Amount    float64 `json:"amount" xml:"amount"`
	UserID    int     `json:"user_id" xml:"user_id"`
	GatewayID int     `json:"gateway_id" xml:"gateway_id"`
	CountryID int     `json:"country_id" xml:"country_id"`
	Currency  string  `json:"currency" xml:"currency"`
}

// APIResponse a standard response structure for the APIs
type APIResponse struct {
	StatusCode int                    `json:"status_code" xml:"status_code"`
	Message    string                 `json:"message" xml:"message"`
	Data       map[string]interface{} `json:"data,omitempty" xml:"data,omitempty"`
}
