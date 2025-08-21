package models

type User struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterUser struct {
	Name string `json:"name" validate:"required,min=5,max=50"`
	User
}
