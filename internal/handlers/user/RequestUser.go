package userHandlers

// UserRequest represents req user request
type UserRequest struct {
	Email    string `json:"email" binding:"required" example:"jonn@gmail.com"`
	Password string `json:"password" binding:"required" example:"securitycod123"`
}

// UserRegistration represents registration user request
type UserRegistration struct {
	Name string `json:"name" binding:"required" example:"jonn"`
	UserRequest
}

// RefreshToken represents getToken user response
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
