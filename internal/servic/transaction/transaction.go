package transaction

import (
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/go-playground/validator/v10"
)

type DatabaseTransactionRepository interface {
	CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
	GetTransaction(TransactionId uint) (domain.TransactionOutput, error)
	UpdateTransaction(transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
	DeleteTransaction(transactionId uint) error
}

type TransactionServer struct {
	d DatabaseTransactionRepository
}

func CreateTransactionServer(d DatabaseTransactionRepository) *TransactionServer {
	return &TransactionServer{
		d: d,
	}
}

func (ts *TransactionServer) CreateTransactionServic(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	if err := validator.New().Struct(&tran); err != nil {
		return 0, fmt.Errorf("error validate: %w", err)
	}

	id, err := ts.d.CreateTransaction(idUser, idCategory, tran)
	if err != nil {
		return 0, fmt.Errorf("error create transaction: %w", err)
	}

	return id, nil
}

func (ts *TransactionServer) ReadTransactionServer(idTransaction uint) (domain.TransactionOutput, error) {
	transaction, err := ts.d.GetTransaction(idTransaction)
	if err != nil {
		return domain.TransactionOutput{}, fmt.Errorf("error get transaction: %w", err)
	}

	return transaction, nil
}

func (ts *TransactionServer) UpdateTransactionServer(idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	if err := validator.New().Struct(&newTransaction); err != nil {
		return domain.TransactionOutput{}, fmt.Errorf("error validate: %w", err)
	}

	transaction, err := ts.d.UpdateTransaction(idTransaction, newTransaction)
	if err != nil {
		return domain.TransactionOutput{}, fmt.Errorf("error update transaction: %w", err)
	}

	return transaction, nil
}

func (ts *TransactionServer) DeleteTransactionServer(idTransaction uint) error {
	err := ts.d.DeleteTransaction(idTransaction)
	if err != nil {
		return fmt.Errorf("error delete transaction: %w", err)
	}
	return nil
}
