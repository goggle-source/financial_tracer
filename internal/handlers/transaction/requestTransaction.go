package transactionHandlers

// RequestCreateTransaction represents registration transaction request
type RequestCreateTransaction struct {
	IdCategory  uint   `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Count       int    `json:"limit" binding:"required"`
	Description string `json:"description"`
}

// RequestIdTransaction represents registration transaction request
type RequestIdTransaction struct {
	IdTransaction uint `json:"transaction_id"`
}

// RequestUpdateTransaction represents registration transaction request
type RequestUpdateTransaction struct {
	IdTransaction uint   `json:"transaction_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Count         int    `json:"limit" binding:"required"`
	Description   string `json:"description"`
}
