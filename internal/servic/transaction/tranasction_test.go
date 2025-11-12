package transaction

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/cash"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionServic(t *testing.T) {
	type test struct {
		name          string
		tran          domain.TransactionInput
		idUser        uint
		idCategory    uint
		idTransaction uint
		repoErr       error
		tranErr       error
		shouldCallDB  bool
		shouldCache   bool
		cacheErr      error
	}

	arrTest := []test{
		{
			name: "success",
			tran: domain.TransactionInput{
				Name:        "траты на еду",
				Count:       1000,
				Description: "потраченно в субботу в ресторане",
			},
			idUser:        123,
			idCategory:    15312,
			idTransaction: 1,
			repoErr:       nil,
			tranErr:       nil,
			shouldCallDB:  true,
			shouldCache:   true,
			cacheErr:      nil,
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
			repoErr:       postgresql.ErrorNotFound,
			tranErr:       ErrNoFound,
			shouldCallDB:  true,
			shouldCache:   false,
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
			repoErr:       errors.New("some db error"),
			tranErr:       ErrDatabase,
			shouldCallDB:  true,
			shouldCache:   false,
		},
		{
			name: "error limit",
			tran: domain.TransactionInput{
				Name:        "лимит по категории",
				Count:       999999,
				Description: "превышение лимита",
			},
			idUser:        2,
			idCategory:    9,
			idTransaction: 0,
			repoErr:       postgresql.ErrorLimit,
			tranErr:       ErrLimit,
			shouldCallDB:  true,
			shouldCache:   false,
		},
		{
			name: "error validate",
			tran: domain.TransactionInput{
				Name:        "",
				Count:       0,
				Description: "",
			},
			idUser:        5,
			idCategory:    10,
			idTransaction: 0,
			repoErr:       nil,
			tranErr:       validator.ValidationErrors{},
			shouldCallDB:  false,
			shouldCache:   false,
		},
	}

	for _, test := range arrTest {
		t.Run(test.name, func(t *testing.T) {
			repoMock := new(DbMock)
			redisMock := new(cash.RedisMock)

			if test.shouldCallDB {
				repoMock.On("CreateTransaction", mock.Anything, test.idUser, test.idCategory, test.tran).
					Return(test.idTransaction, test.repoErr)
			}
			if test.shouldCache {
				expectedTransaction := domain.TransactionOutput{
					Name:        test.tran.Name,
					UserID:      test.idUser,
					CategoryID:  test.idCategory,
					Description: test.tran.Description,
					Count:       test.tran.Count,
				}
				redisMock.On("HsetTransaction", mock.Anything, test.idTransaction, expectedTransaction).
					Return(test.cacheErr)
			}
			log := logrus.New()

			server := CreateTransactionServer(repoMock, repoMock, repoMock, repoMock, log, redisMock)
			id, err := server.CreateTransaction(context.Background(), test.idUser, test.idCategory, test.tran)

			if test.repoErr != nil || test.tranErr != nil {
				assert.Error(t, err)
				if test.name == "error validate" {
					var verr validator.ValidationErrors
					if !errors.As(err, &verr) {
						t.Fatalf("err != test.tranErr: %v", err)
					}
				} else if !errors.Is(err, test.tranErr) {
					t.Fatalf("err != test.tranErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id, test.idTransaction)
			}

			if test.shouldCallDB {
				repoMock.AssertExpectations(t)
			} else {
				repoMock.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}

			if test.shouldCache {
				redisMock.AssertExpectations(t)
			} else {
				redisMock.AssertNotCalled(t, "HsetTransaction", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestReadTransactionServer(t *testing.T) {
	type test struct {
		name          string
		tran          domain.TransactionOutput
		idTransaction uint
		tranErr       error
		svcErr        error
		redisPayload  map[string]string
		redisErr      error
		shouldCallDB  bool
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
			idTransaction: 7,
			tranErr:       nil,
			svcErr:        nil,
			redisPayload: map[string]string{
				"name":        "траты на магазин",
				"description": "походил с девушкой по магазинам",
				"userID":      strconv.FormatUint(uint64(1), 10),
				"categoryID":  strconv.FormatUint(uint64(2), 10),
				"count":       strconv.Itoa(10000),
			},
			redisErr:     nil,
			shouldCallDB: false,
		},
		{
			name: "error database",
			tran: domain.TransactionOutput{
				UserID:      5,
				CategoryID:  1,
				Name:        "покупка нового пк",
				Count:       100000,
				Description: "купил себе компьютер по-мощнее для разработки собственной нейросети",
			},
			idTransaction: 0,
			tranErr:       errors.New("some db err"),
			svcErr:        ErrDatabase,
			redisPayload:  map[string]string{},
			redisErr:      errors.New("cache miss"),
			shouldCallDB:  true,
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
			tranErr:       postgresql.ErrorNotFound,
			svcErr:        ErrNoFound,
			redisPayload:  map[string]string{},
			redisErr:      errors.New("cache miss"),
			shouldCallDB:  true,
		},
	}

	for _, ts := range arrTest {
		t.Run(ts.name, func(t *testing.T) {

			repoMock := new(DbMock)
			redisMock := new(cash.RedisMock)

			if ts.shouldCallDB {
				repoMock.On("GetTransaction", mock.Anything, ts.idTransaction).Return(ts.tran, ts.tranErr)
			}
			redisMock.On("HgetTransaction", mock.Anything, ts.idTransaction).
				Return(ts.redisPayload, ts.redisErr)
			log := logrus.New()

			server := CreateTransactionServer(repoMock, repoMock, repoMock, repoMock, log, redisMock)
			tran, err := server.GetTransaction(context.Background(), ts.idTransaction)
			if ts.tranErr != nil || ts.svcErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, ts.svcErr) {
					t.Fatalf("err != ts.tranErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tran, ts.tran)
			}

			redisMock.AssertExpectations(t)

			if ts.shouldCallDB {
				repoMock.AssertExpectations(t)
			} else {
				repoMock.AssertNotCalled(t, "GetTransaction", mock.Anything, mock.AnythingOfType("uint"))
			}
		})
	}
}

func TestUpdateTransactionServer(t *testing.T) {
	type test struct {
		name          string
		idTransaction uint
		tranInput     domain.TransactionInput
		tranOutput    domain.TransactionOutput
		tranErr       error
		svcErr        error
		shouldCallDB  bool
		shouldCache   bool
		cacheErr      error
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
			tranErr:      nil,
			svcErr:       nil,
			shouldCallDB: true,
			shouldCache:  true,
			cacheErr:     nil,
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
			tranErr:      errors.New("db error"),
			svcErr:       ErrDatabase,
			shouldCallDB: true,
			shouldCache:  false,
		},
		{
			name:          "error validate",
			idTransaction: 0,
			tranInput: domain.TransactionInput{
				Name:        "",
				Count:       0,
				Description: "",
			},
			tranOutput:   domain.TransactionOutput{},
			tranErr:      nil,
			svcErr:       validator.ValidationErrors{},
			shouldCallDB: false,
			shouldCache:  false,
		},
	}

	for _, test := range arrTest {
		t.Run(test.name, func(t *testing.T) {

			repoMock := new(DbMock)
			redisMock := new(cash.RedisMock)

			if test.shouldCallDB {
				repoMock.On("UpdateTransaction", mock.Anything, test.idTransaction, test.tranInput).Return(test.tranOutput, test.tranErr)
			}
			if test.shouldCache {
				redisMock.On("HsetTransaction", mock.Anything, test.idTransaction, test.tranOutput).Return(test.cacheErr)
			}
			log := logrus.New()

			server := CreateTransactionServer(repoMock, repoMock, repoMock, repoMock, log, redisMock)
			tranOutput, err := server.UpdateTransaction(context.Background(), test.idTransaction, test.tranInput)

			if test.tranErr != nil || test.svcErr != nil {
				assert.Error(t, err)
				if test.name == "error validate" {
					var verr validator.ValidationErrors
					if !errors.As(err, &verr) {
						t.Fatalf("err != validator.ValidationErrors: %v", err)
					}
				} else if !errors.Is(err, test.svcErr) {
					t.Fatalf("err != test.tranErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tranOutput, test.tranOutput)
			}

			if test.shouldCallDB {
				repoMock.AssertExpectations(t)
			} else {
				repoMock.AssertNotCalled(t, "UpdateTransaction", mock.Anything, mock.Anything, mock.Anything)
			}

			if test.shouldCache {
				redisMock.AssertExpectations(t)
			} else {
				redisMock.AssertNotCalled(t, "HsetTransaction", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}

func TestDeleteTransactionServer(t *testing.T) {
	type test struct {
		name          string
		idTransaction uint
		tranErr       error
		svcErr        error
		shouldCache   bool
		cacheErr      error
	}

	arrTest := []test{
		{
			name:          "success",
			idTransaction: 2,
			tranErr:       nil,
			svcErr:        nil,
			shouldCache:   true,
			cacheErr:      nil,
		},
		{
			name:          "not found",
			idTransaction: 5,
			tranErr:       postgresql.ErrorNotFound,
			svcErr:        ErrNoFound,
			shouldCache:   false,
		},
		{
			name:          "error database",
			idTransaction: 6,
			tranErr:       errors.New("db error"),
			svcErr:        ErrDatabase,
			shouldCache:   false,
		},
	}

	for _, ts := range arrTest {
		t.Run(ts.name, func(t *testing.T) {

			repoMock := new(DbMock)
			redisMock := new(cash.RedisMock)

			repoMock.On("DeleteTransaction", mock.Anything, ts.idTransaction).Return(ts.tranErr)
			if ts.shouldCache {
				redisMock.On("HdelTransaction", mock.Anything, ts.idTransaction).Return(ts.cacheErr)
			}
			log := logrus.New()

			server := CreateTransactionServer(repoMock, repoMock, repoMock, repoMock, log, redisMock)
			err := server.DeleteTransaction(context.Background(), ts.idTransaction)
			if ts.tranErr != nil || ts.svcErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, ts.svcErr) {
					t.Fatalf("err != ts.tranErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
			}
			repoMock.AssertExpectations(t)
			if ts.shouldCache {
				redisMock.AssertExpectations(t)
			} else {
				redisMock.AssertNotCalled(t, "HdelTransaction", mock.Anything, mock.Anything)
			}
		})
	}
}
