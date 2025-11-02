package transactionHandlers

// RequestCreateTransaction represents registration transaction request
type RequestCreateTransaction struct {
	IdCategory  uint   `json:"category_id" binding:"required" example:"2"`
	Name        string `json:"name" binding:"required" example:"jonn"`
	Count       int    `json:"limit" binding:"required" example:"1000"`
	Description string `json:"description" example:"spending on food"`
}

// RequestUpdateTransaction represents registration transaction request
type RequestUpdateTransaction struct {
	IdTransaction uint   `json:"transaction_id" binding:"required" example:"1"`
	Name          string `json:"name" binding:"required" example:"jonn"`
	Count         int    `json:"limit" binding:"required" example:"2000"`
	Description   string `json:"description" example:"going to a restaurant"`
}
