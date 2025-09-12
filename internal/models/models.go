package models

type AuthenticationUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterUser struct {
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type DeleteUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Category represents a category model
// @Name Category
type Category struct {
	Name        string `json:"name" validate:"required,max=60,min=3"`
	Limit       int    `json:"limit" validate:"required"`
	Description string `json:"description" validate:"required,max=100"`
}
