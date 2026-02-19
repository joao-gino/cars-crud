package domain

type ValidateRequest struct {
	APIKey string `json:"api_key" example:"my-api-key-12345"`
}

type ValidateResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}
