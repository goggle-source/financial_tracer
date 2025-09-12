package userHandlers

// UserRegistration represents registration user request
type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserAuthentication represents authentication user request
type UserAuthentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserDelete represents delete user request
type UserDelete struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshToken represents getToken user request
type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
