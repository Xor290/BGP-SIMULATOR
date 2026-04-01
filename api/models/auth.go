package models

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterClientRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=8"`
	Nom       string `json:"nom" binding:"required,min=2,max=100"`
	Prenom    string `json:"prenom" binding:"required,min=2,max=100"`
	Telephone string `json:"telephone" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	User        any    `json:"user"`
}

type ProfileResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
