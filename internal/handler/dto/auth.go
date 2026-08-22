package dto

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn string `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
}
