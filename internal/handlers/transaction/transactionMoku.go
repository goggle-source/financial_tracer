package transactionHandlers

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type tranasctionServicMock struct {
	mock.Mock
}

func (d *tranasctionServicMock) CreateTransaction(ctx context.Context, idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	args := d.Called(ctx, idUser, idCategory, tran)
	return args.Get(0).(uint), args.Error(1)
}
func (d *tranasctionServicMock) GetTransaction(ctx context.Context, idTransaction uint) (domain.TransactionOutput, error) {
	args := d.Called(ctx, idTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}
func (d *tranasctionServicMock) UpdateTransaction(ctx context.Context, idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	args := d.Called(ctx, idTransaction, newTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}
func (d *tranasctionServicMock) DeleteTransaction(ctx context.Context, idTransaction uint) error {
	args := d.Called(ctx, idTransaction)
	return args.Error(0)
}
