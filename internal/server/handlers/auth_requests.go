package handlers

type RegisterRequest struct {
	Login            string `json:"login" validate:"required"`
	Password         string `json:"password" validate:"required"`
	TwilioAccountSID string `json:"twilio_account_sid" validate:"required"`
	TwilioAuthToken  string `json:"twilio_auth_token" validate:"required"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
