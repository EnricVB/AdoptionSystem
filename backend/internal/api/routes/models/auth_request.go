package r_models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleLoginRequest struct {
	Email   string `json:"email"`
	IDToken string `json:"id_token"`
}

type TwoFactorRequest struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}

type RefreshTokenRequest struct {
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	Address string `json:"address"`

	Password   string `json:"password"`
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}
