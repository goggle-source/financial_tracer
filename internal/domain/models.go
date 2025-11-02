package domain

type AuthenticationUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type RegisterUser struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type DeleteUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type User struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hasy"`
}

// @Name	Category
type CategoryInput struct {
	Name        string `json:"name" validate:"required,max=60,min=3"`
	Limit       int    `json:"limit" validate:"required"`
	Type        string `json:"type" validate:"max=100"`
	Description string `json:"description" validate:"max=100"`
}

type CategoryOutput struct {
	UserID      uint
	Name        string `json:"name"`
	Limit       int    `json:"limit"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// @Name	Transaction
type TransactionInput struct {
	Name        string `json:"name" validate:"required,max=60,min=3"`
	Count       int    `json:"count" validate:"required"`
	Description string `json:"description" validate:"max=100"`
}

type TransactionOutput struct {
	UserID      uint
	CategoryID  uint
	Name        string `json:"name" validate:"required,max=60,min=3"`
	Count       int    `json:"count" validate:"required"`
	Description string `json:"description" validate:"max=100"`
}
