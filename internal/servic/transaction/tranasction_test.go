package transaction

import (
	"errors"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransactionServic(t *testing.T) {
	type test struct {
		name          string
		tran          domain.TransactionInput
		idUser        uint
		idCategory    uint
		idTransaction uint
		mockErr       error
		msgErr        string
		shouldCallDB  bool
	}

	arrTest := []test{
		{
			name: "success",
			tran: domain.TransactionInput{
				Name:        "траты на еду",
				Count:       1000,
				Description: "потраченно в субботу в ресторане",
			},
			idUser:        1,
			idCategory:    2,
			idTransaction: 1,
			mockErr:       nil,
			msgErr:        "",
			shouldCallDB:  true,
		},
		{
			name: "error not found",
			tran: domain.TransactionInput{
				Name:        "покупка курсов по Python",
				Count:       100000,
				Description: "чтобы стать senior developer on python",
			},
			idUser:        1,
			idCategory:    2,
			idTransaction: 0,
			mockErr:       errors.New("error not found"),
			msgErr:        "error create transaction",
			shouldCallDB:  true,
		},
		{
			name: "error database",
			tran: domain.TransactionInput{
				Name:        "траты на собаку",
				Count:       1000,
				Description: "Куплены игрушки для собаки",
			},
			idUser:        4,
			idCategory:    6,
			idTransaction: 0,
			mockErr:       errors.New("error database"),
			msgErr:        "error create transaction",
			shouldCallDB:  true,
		},
		{
			name: "validate",
			tran: domain.TransactionInput{
				Name:        "",
				Count:       0,
				Description: "",
			},
			idUser:        5,
			idCategory:    10,
			idTransaction: 0,
			mockErr:       nil,
			msgErr:        "error validate",
			shouldCallDB:  false,
		},
	}

	for _, test := range arrTest {
		t.Run(test.name, func(t *testing.T) {
			repoMock := new(DbMock)

			repoMock.On("CreateTransaction", test.idUser, test.idCategory, test.tran).
				Return(test.idTransaction, test.mockErr)

			server := CreateTransactionServer(repoMock)
			id, err := server.CreateTransactionServic(test.idUser, test.idCategory, test.tran)

			if test.mockErr != nil || test.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.msgErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id, test.idTransaction)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "CreateTransaction", test.idUser, test.idCategory, test.tran)
			}
		})
	}
}

func TestReadTransactionServer(t *testing.T) {
	type test struct {
		name          string
		tran          domain.TransactionOutput
		idTransaction uint
		mockErr       error
		msgErr        string
	}

	arrTest := []test{
		{
			name: "success",
			tran: domain.TransactionOutput{
				UserID:      1,
				CategoryID:  2,
				Name:        "траты на магазин",
				Count:       10000,
				Description: "походил с девушкой по магазинам",
			},
			idTransaction: 4,
			mockErr:       nil,
			msgErr:        "",
		},
		{
			name: "error database",
			tran: domain.TransactionOutput{
				UserID:      4,
				CategoryID:  7,
				Name:        "покупка нового пк",
				Count:       100000,
				Description: "купил себе компьютер по-мощнее для разработки собственной нейросети",
			},
			idTransaction: 0,
			mockErr:       errors.New("error database"),
			msgErr:        "error get transaction",
		},
		{
			name: "not found",
			tran: domain.TransactionOutput{
				UserID:      4,
				CategoryID:  7,
				Name:        "покупка нового пк",
				Count:       1000000,
				Description: "купил себе компьютер для игр",
			},
			idTransaction: 0,
			mockErr:       errors.New("error not found"),
			msgErr:        "error get transaction",
		},
	}

	for _, ts := range arrTest {
		t.Run(ts.name, func(t *testing.T) {
			repoMock := new(DbMock)

			repoMock.On("GetTransaction", ts.idTransaction).Return(ts.tran, ts.mockErr)

			server := CreateTransactionServer(repoMock)
			tran, err := server.ReadTransactionServer(ts.idTransaction)
			if ts.mockErr != nil || ts.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), ts.msgErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tran, ts.tran)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestUpdateTransactionServer(t *testing.T) {
	type test struct {
		name          string
		idTransaction uint
		tranInput     domain.TransactionInput
		tranOutput    domain.TransactionOutput
		mockErr       error
		msgErr        string
		shouldCallDB  bool
	}

	arrTest := []test{
		{
			name:          "success",
			idTransaction: 2,
			tranInput: domain.TransactionInput{
				Name:        "траты на еду",
				Count:       1000,
				Description: "сходил в шаурмичную",
			},
			tranOutput: domain.TransactionOutput{
				UserID:      2,
				CategoryID:  6,
				Name:        "траты на еду",
				Count:       2000,
				Description: "сходил в шаурмичную 2 раза",
			},
			mockErr:      nil,
			msgErr:       "",
			shouldCallDB: true,
		},
		{
			name:          "error database",
			idTransaction: 0,
			tranInput: domain.TransactionInput{
				Name:        "продукты",
				Count:       2500,
				Description: "сходил в машазин",
			},
			tranOutput:   domain.TransactionOutput{},
			mockErr:      errors.New("error database"),
			msgErr:       "error update transaction",
			shouldCallDB: true,
		},
		{
			name:          "validate",
			idTransaction: 0,
			tranInput: domain.TransactionInput{
				Name:        "",
				Count:       0,
				Description: "",
			},
			tranOutput:   domain.TransactionOutput{},
			mockErr:      nil,
			msgErr:       "error validate",
			shouldCallDB: false,
		},
	}

	for _, test := range arrTest {
		t.Run(test.name, func(t *testing.T) {
			repoMock := new(DbMock)

			repoMock.On("UpdateTransaction", test.idTransaction, test.tranInput).Return(test.tranOutput, test.mockErr)

			server := CreateTransactionServer(repoMock)
			tranOutput, err := server.UpdateTransactionServer(test.idTransaction, test.tranInput)

			if test.mockErr != nil || test.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.msgErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tranOutput, test.tranOutput)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "UpdateTransaction", test.idTransaction, test.tranInput)
			}
		})
	}
}

func TestDeleteTransactionServer(t *testing.T) {
	type test struct {
		name          string
		idTransaction uint
		mockErr       error
		msgErr        string
	}

	arrTest := []test{
		{
			name:          "success",
			idTransaction: 2,
			mockErr:       nil,
			msgErr:        "",
		},
		{
			name:          "not found",
			idTransaction: 5,
			mockErr:       errors.New("error not found"),
			msgErr:        "error delete transaction",
		},
		{
			name:          "error database",
			idTransaction: 6,
			mockErr:       errors.New("error database"),
			msgErr:        "error delete transaction",
		},
	}

	for _, ts := range arrTest {
		t.Run(ts.name, func(t *testing.T) {
			repoMock := new(DbMock)

			repoMock.On("DeleteTransaction", ts.idTransaction).Return(ts.mockErr)

			server := CreateTransactionServer(repoMock)
			err := server.DeleteTransactionServer(ts.idTransaction)
			if ts.mockErr != nil || ts.msgErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), ts.msgErr)
			} else {
				assert.NoError(t, err)
			}
			repoMock.AssertExpectations(t)
		})
	}
}
