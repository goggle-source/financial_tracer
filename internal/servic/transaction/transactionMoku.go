package transaction

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	args := d.Called(idUser, idCategory, tran)
	return args.Get(0).(uint), args.Error(1)
}

func (d *DbMock) GetTransaction(TransactionId uint) (domain.TransactionOutput, error) {
	args := d.Called(TransactionId)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}

func (d *DbMock) UpdateTransaction(transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	args := d.Called(transactionId, newTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}

func (d *DbMock) DeleteTransaction(transactionId uint) error {
	args := d.Called(transactionId)
	return args.Error(0)
}
