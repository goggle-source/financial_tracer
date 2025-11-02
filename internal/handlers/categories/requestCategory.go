package categoryHandlers

// RequestCreateCategory represents CreateCategory category request
type RequestCreateCategory struct {
	Name        string `json:"name" binding:"required" example:"jonn"`
	Limit       int    `json:"limit" binding:"required" exemple:"2000"`
	Description string `json:"description" binding:"required" example:"Shopping in the store"`
}

// ResponseUpdateCategory represents UpdateCategory
type RequestUpdateCategory struct {
	Name        string `json:"name" binding:"required" example:"jonn"`
	Limit       int    `json:"limit" binding:"required" example:"2000"`
	Description string `json:"description" binding:"required" example:"car expenses"`
	CategoryId  uint   `json:"category_id" example:"2"`
}
