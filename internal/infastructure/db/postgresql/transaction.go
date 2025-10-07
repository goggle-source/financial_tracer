package postgresql

import (
	"errors"

	"github.com/financial_tracer/internal/domain"
	"gorm.io/gorm"
)

func (d *Db) CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {

	tx := d.DB.Begin()
	newTransaction := Transaction{
		UserID:      idUser,
		CategoryID:  idCategory,
		Name:        tran.Name,
		Count:       tran.Count,
		Description: tran.Description,
	}
	var categor Category
	result := tx.Select("limit").First(categor, idCategory)
	if result.Error != nil {
		tx.Rollback()
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, ErrorNotFound
		}
		return 0, result.Error
	}
	if categor.Limit < newTransaction.Count {
		return 0, ErrorLimit
	}

	result = tx.Create(&newTransaction)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	return newTransaction.ID, tx.Commit().Error

}

func (d *Db) GetTransaction(TransactionId uint) (domain.TransactionOutput, error) {
	var tran Transaction

	result := d.DB.First(&tran, TransactionId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.TransactionOutput{}, ErrorNotFound
		}
		return domain.TransactionOutput{}, result.Error
	}

	transaction := domain.TransactionOutput{
		UserID:      tran.UserID,
		CategoryID:  tran.CategoryID,
		Name:        tran.Name,
		Count:       tran.Count,
		Description: tran.Description,
	}
	return transaction, nil
}

func (d *Db) UpdateTransaction(transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	var transaction = Transaction{
		Name:        newTransaction.Name,
		Count:       newTransaction.Count,
		Description: newTransaction.Description,
	}

	result := d.DB.Updates(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.TransactionOutput{}, ErrorNotFound
		}
		return domain.TransactionOutput{}, result.Error
	}
	res := domain.TransactionOutput{
		UserID:      transaction.UserID,
		CategoryID:  transaction.CategoryID,
		Name:        transaction.Name,
		Count:       transaction.Count,
		Description: transaction.Description,
	}

	return res, nil
}

func (d *Db) DeleteTransaction(transactionId uint) error {

	result := d.DB.Delete(&Transaction{}, transactionId)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return ErrorNotFound
		}
		return result.Error
	}

	return nil
}
