package categoryHandlers

// RequestCreateCategory represents CreateCategory category request
type RequestCreateCategory struct {
	Name        string `json:"name" binding:"required"`
	Limit       int    `json:"limit" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// IDCategory represents GetIdCategory category request
type IDCategory struct {
	CategoryId uint `json:"category_id"`
}

// ResponseUpdateCategory represents UpdateCategory
type RequestUpdateCategory struct {
	Name        string `json:"name" binding:"required"`
	Limit       int    `json:"limit" binding:"required"`
	Description string `json:"description" binding:"required"`
	CategoryId  uint   `json:"category_id"`
}
