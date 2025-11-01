package transactionHandlers

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type tranasctionServicMock struct {
	mock.Mock
}

func (d *tranasctionServicMock) CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	args := d.Called(idUser, idCategory, tran)
	return args.Get(0).(uint), args.Error(1)
}
func (d *tranasctionServicMock) GetTransaction(idTransaction uint) (domain.TransactionOutput, error) {
	args := d.Called(idTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}
func (d *tranasctionServicMock) UpdateTransaction(idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	args := d.Called(idTransaction, newTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}
func (d *tranasctionServicMock) DeleteTransaction(idTransaction uint) error {
	args := d.Called(idTransaction)
	return args.Error(0)
}
