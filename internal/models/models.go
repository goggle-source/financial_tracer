package models

type UserRequest struct {
	Name     string `json:"name" binding:"required,max=50,min=0"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
