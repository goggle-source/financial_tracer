package userHandlers

type UserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserRegistration represents registration user request
type UserRegistration struct {
	Name string `json:"name" binding:"required"`
	UserRequest
}

// RefreshToken represents getToken user request
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
