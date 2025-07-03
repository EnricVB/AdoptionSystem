package r_models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TwoFactorRequest struct {
	Email string `json:"email"`
	_2fa  string `json:"two_factor_auth"`
}
