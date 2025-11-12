package transaction

import (
	"context"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func (d *DbMock) CreateTransaction(ctx context.Context, idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	args := d.Called(ctx, idUser, idCategory, tran)
	return args.Get(0).(uint), args.Error(1)
}

func (d *DbMock) GetTransaction(ctx context.Context, TransactionId uint) (domain.TransactionOutput, error) {
	args := d.Called(ctx, TransactionId)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}

func (d *DbMock) UpdateTransaction(ctx context.Context, transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	args := d.Called(ctx, transactionId, newTransaction)
	return args.Get(0).(domain.TransactionOutput), args.Error(1)
}

func (d *DbMock) DeleteTransaction(ctx context.Context, transactionId uint) error {
	args := d.Called(ctx, transactionId)
	return args.Error(0)
}
