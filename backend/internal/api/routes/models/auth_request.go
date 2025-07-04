package r_models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TwoFactorRequest struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}

type RefreshTokenRequest struct {
	Email string `json:"email"`
}
