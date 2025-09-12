package categoryHandlers

// RequestCreateCategory represents CreateCategory category request
type RequestCreateCategory struct {
	Name        string `json:"name"`
	Limit       int    `json:"limit"`
	Description string `json:"description"`
}

// IDCategory represents GetIdCategory category request
type IDCategory struct {
	CategoryId uint `json:"category_id"`
}

// ResponseUpdateCategory represents UpdateCategory
type RequestUpdateCategory struct {
	Name        string `json:"name"`
	Limit       int    `json:"limit"`
	Description string `json:"Description"`
	CategoryId  uint   `json:"category_id"`
}
