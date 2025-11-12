package category

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

func TestCreateCategory(t *testing.T) {

	type tests struct {
		Name            string
		userID          uint
		category        domain.CategoryInput
		categoryID      uint
		mockErr         error
		categoryErr     error
		shouldCallDB    bool
		RedisErr        error
		shouldCallRedis bool
	}

	arrTests := []tests{
		{
			Name:   "success",
			userID: 1,
			category: domain.CategoryInput{
				Name:        "chicken",
				Limit:       10000,
				Type:        "траты на еду",
				Description: "сходил в ресторан",
			},
			categoryID:      2,
			mockErr:         nil,
			categoryErr:     nil,
			RedisErr:        nil,
			shouldCallDB:    true,
			shouldCallRedis: true,
		},
		{
			Name:   "success with redis error (should not fail)",
			userID: 1,
			category: domain.CategoryInput{
				Name:        "chicken",
				Limit:       10000,
				Type:        "траты на еду",
				Description: "сходил в ресторан",
			},
			categoryID:      2,
			mockErr:         nil,
			categoryErr:     nil,
			RedisErr:        errors.New("redis connection error"),
			shouldCallDB:    true,
			shouldCallRedis: true,
		},
		{
			Name:   "not found",
			userID: 2,
			category: domain.CategoryInput{
				Name:        "default",
				Limit:       10000,
				Description: "test case",
			},
			categoryID:      0,
			mockErr:         postgresql.ErrorNotFound,
			categoryErr:     ErrNoFound,
			RedisErr:        nil,
			shouldCallDB:    true,
			shouldCallRedis: false,
		},
		{
			Name:   "duplicated",
			userID: 3,
			category: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			categoryID:      0,
			mockErr:         postgresql.ErrorDuplicated,
			categoryErr:     ErrDuplicated,
			RedisErr:        nil,
			shouldCallDB:    true,
			shouldCallRedis: false,
		},
		{
			Name:   "valdiate",
			userID: 5,
			category: domain.CategoryInput{
				Name:        "a",
				Limit:       100000,
				Description: "asfdasdfgAFG",
			},
			categoryID:      0,
			mockErr:         nil,
			categoryErr:     validator.ValidationErrors{},
			RedisErr:        nil,
			shouldCallDB:    false,
			shouldCallRedis: false,
		},
		{
			Name:   "error database",
			userID: 10,
			category: domain.CategoryInput{
				Name:        "sosalka",
				Limit:       2000,
				Description: "what my write?",
			},
			categoryID:      0,
			mockErr:         errors.New("error database"),
			categoryErr:     ErrDatabase,
			RedisErr:        nil,
			shouldCallDB:    true,
			shouldCallRedis: false,
		},
	}
	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			r := new(cash.RedisMock)

			repoMock.ExpectedCalls = nil
			repoMock.Calls = nil
			r.ExpectedCalls = nil
			r.Calls = nil
			category := domain.CategoryOutput{
				Name:        test.category.Name,
				Description: test.category.Description,
				Type:        test.category.Type,
				Limit:       test.category.Limit,
				UserID:      test.userID,
			}

			if test.shouldCallDB {
				repoMock.On("CreateCategory", mock.Anything, test.userID, test.category).Return(test.categoryID, test.mockErr)
			}

			if test.shouldCallRedis {
				r.On("HsetCategory", context.Background(), test.categoryID, category).Return(test.RedisErr)
			}

			log := logrus.New()

			servic := CreateCategoryServer(repoMock, repoMock, repoMock, repoMock, repoMock, log, r)
			resultID, err := servic.CreateCategory(context.Background(), test.userID, test.category)

			if test.mockErr != nil || test.categoryErr != nil {
				assert.Error(t, err)
				if test.Name == "valdiate" {
					var verr validator.ValidationErrors
					if !errors.As(err, &verr) {
						t.Fatalf("err != validator.ValidationErrors: %v", err)
					}
				} else if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.categoryID, resultID)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "CreateCategory", mock.Anything, test.userID, test.category)
			} else {
				repoMock.AssertNotCalled(t, "CreateCategory", mock.Anything, test.userID, test.category)
			}

			if test.shouldCallRedis {
				r.AssertCalled(t, "HsetCategory", context.Background(), test.categoryID, category)
			}

			repoMock.AssertExpectations(t)
			r.AssertExpectations(t)
		})
	}
}

func TestGetCategory(t *testing.T) {

	type tests struct {
		Name            string
		category        domain.CategoryOutput
		categoryID      uint
		mockErr         error
		categoryErr     error
		redisData       map[string]string
		redisErr        error
		shouldCallRedis bool
		shouldCallDB    bool
	}

	arrTests := []tests{
		{
			Name: "success from cache",
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
				Type:        "траты на еду",
			},
			categoryID:  2,
			mockErr:     nil,
			categoryErr: nil,
			redisData: map[string]string{
				"userID":      "2",
				"name":        "food",
				"limit":       "10000",
				"description": "траты на еду",
				"type":        "траты на еду",
			},
			redisErr:        nil,
			shouldCallRedis: true,
			shouldCallDB:    false,
		},
		{
			Name: "success from database (cache miss)",
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			categoryID:      6,
			mockErr:         nil,
			categoryErr:     nil,
			redisData:       nil,
			redisErr:        errors.New("cache miss"),
			shouldCallRedis: true,
			shouldCallDB:    true,
		},
		{
			Name:            "not found",
			category:        domain.CategoryOutput{},
			categoryID:      0,
			mockErr:         postgresql.ErrorNotFound,
			categoryErr:     ErrNoFound,
			redisData:       nil,
			redisErr:        errors.New("cache miss"),
			shouldCallRedis: true,
			shouldCallDB:    true,
		},
		{
			Name: "error database",
			category: domain.CategoryOutput{
				UserID:      1,
				Name:        "gector",
				Description: "mult in cinema",
				Limit:       2000,
				Type:        "happy",
			},
			categoryID:      7,
			mockErr:         errors.New("error database"),
			categoryErr:     ErrDatabase,
			redisData:       nil,
			redisErr:        errors.New("cache miss"),
			shouldCallRedis: true,
			shouldCallDB:    true,
		},
	}

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			r := new(cash.RedisMock)

			repoMock.ExpectedCalls = nil
			repoMock.Calls = nil

			if test.shouldCallDB {
				repoMock.On("GetCategory", mock.Anything, test.categoryID).Return(test.category, test.mockErr)
			}

			log := logrus.New()

			if test.shouldCallRedis {
				r.On("HgetCategory", context.Background(), test.categoryID).Return(test.redisData, test.redisErr)
			}

			servic := CreateCategoryServer(repoMock, repoMock, repoMock, repoMock, repoMock, log, r)
			resultCategory, err := servic.GetCategory(context.Background(), test.categoryID)

			if test.categoryErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}

			} else {
				assert.NoError(t, err)
				if test.shouldCallRedis && test.redisErr == nil {
					userID, _ := strconv.Atoi(test.redisData["userID"])
					assert.Equal(t, uint(userID), resultCategory.UserID)
					assert.Equal(t, test.redisData["name"], resultCategory.Name)
					assert.Equal(t, test.redisData["description"], resultCategory.Description)
					assert.Equal(t, test.redisData["type"], resultCategory.Type)
					limit, _ := strconv.Atoi(test.redisData["limit"])
					assert.Equal(t, limit, resultCategory.Limit)
				} else {
					assert.Equal(t, test.category, resultCategory)
				}
			}

			if test.shouldCallRedis {
				r.AssertCalled(t, "HgetCategory", context.Background(), test.categoryID)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "GetCategory", mock.Anything, test.categoryID)
			} else {
				repoMock.AssertNotCalled(t, "GetCategory", mock.Anything, test.categoryID)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {

	type tests struct {
		Name            string
		newCategory     domain.CategoryInput
		category        domain.CategoryOutput
		categoryID      uint
		mockErr         error
		categoryErr     error
		shouldCallDB    bool
		shouldCallRedis bool
		redisErr        error
	}

	arrTests := []tests{
		{
			Name: "success",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Type:        "shope",
				Description: "траты на еду",
			},
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       15000,
				Description: "на еду и средства",
			},
			categoryID:      2,
			mockErr:         nil,
			categoryErr:     nil,
			shouldCallDB:    true,
			shouldCallRedis: true,
			redisErr:        nil,
		},
		{
			Name: "success with redis error (should not fail)",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Type:        "shope",
				Description: "траты на еду",
			},
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       15000,
				Description: "на еду и средства",
			},
			categoryID:      2,
			mockErr:         nil,
			categoryErr:     nil,
			shouldCallDB:    true,
			shouldCallRedis: true,
			redisErr:        errors.New("redis connection error"),
		},
		{
			Name: "not found",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Type:        "shope",
				Description: "ssssss",
			},
			category:        domain.CategoryOutput{},
			categoryID:      0,
			mockErr:         postgresql.ErrorNotFound,
			categoryErr:     ErrNoFound,
			shouldCallDB:    true,
			shouldCallRedis: false,
			redisErr:        nil,
		},
		{
			Name: "duplicated",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "",
			},
			category:        domain.CategoryOutput{},
			categoryID:      1,
			mockErr:         postgresql.ErrorDuplicated,
			categoryErr:     ErrDuplicated,
			shouldCallDB:    true,
			shouldCallRedis: false,
			redisErr:        nil,
		},
		{
			Name: "valdiate",
			newCategory: domain.CategoryInput{
				Name:        "",
				Limit:       0,
				Description: "",
			},
			category:        domain.CategoryOutput{},
			categoryID:      0,
			mockErr:         nil,
			categoryErr:     validator.ValidationErrors{},
			shouldCallDB:    false,
			shouldCallRedis: false,
			redisErr:        nil,
		},
	}

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			r := new(cash.RedisMock)

			repoMock.ExpectedCalls = nil
			repoMock.Calls = nil

			if test.shouldCallDB {
				repoMock.On("UpdateCategory", mock.Anything, test.categoryID, test.newCategory).Return(test.category, test.mockErr)
			}

			log := logrus.New()

			if test.shouldCallRedis {
				r.On("HsetCategory", context.Background(), test.categoryID, test.category).Return(test.redisErr)
			}

			servic := CreateCategoryServer(repoMock, repoMock, repoMock, repoMock, repoMock, log, r)
			_, err := servic.UpdateCategory(context.Background(), test.categoryID, test.newCategory)

			if test.mockErr != nil || test.categoryErr != nil {
				assert.Error(t, err)
				if test.Name == "valdiate" {
					var verr validator.ValidationErrors
					if !errors.As(err, &verr) {
						t.Fatalf("err != validator.ValidationErrors: %v", err)
					}
				} else if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "UpdateCategory", mock.Anything, test.categoryID, test.newCategory)
			}

			if test.shouldCallRedis {
				r.AssertCalled(t, "HsetCategory", context.Background(), test.categoryID, test.category)
			}
		})
	}
}

func TestDeleteCategoryDatabase(t *testing.T) {
	type tests struct {
		Name            string
		categoryID      uint
		mockErr         error
		categoryErr     error
		shouldCallRedis bool
		redisErr        error
	}

	arrTests := []tests{
		{
			Name:            "success",
			categoryID:      2,
			mockErr:         nil,
			categoryErr:     nil,
			shouldCallRedis: true,
			redisErr:        nil,
		},
		{
			Name:            "success with redis error (should not fail)",
			categoryID:      3,
			mockErr:         nil,
			categoryErr:     nil,
			shouldCallRedis: true,
			redisErr:        errors.New("redis connection error"),
		},
		{
			Name:            "not found",
			categoryID:      0,
			mockErr:         postgresql.ErrorNotFound,
			categoryErr:     ErrNoFound,
			shouldCallRedis: false,
			redisErr:        nil,
		},
	}

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			r := new(cash.RedisMock)

			repoMock.ExpectedCalls = nil
			repoMock.Calls = nil

			repoMock.On("DeleteCategory", mock.Anything, test.categoryID).Return(test.mockErr)
			log := logrus.New()

			if test.shouldCallRedis {
				r.On("HdelCategory", context.Background(), test.categoryID).Return(test.redisErr)
			}

			servic := CreateCategoryServer(repoMock, repoMock, repoMock, repoMock, repoMock, log, r)
			err := servic.DeleteCategory(context.Background(), test.categoryID)

			if test.mockErr != nil || test.categoryErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
			}
			repoMock.AssertCalled(t, "DeleteCategory", mock.Anything, test.categoryID)

			if test.shouldCallRedis {
				r.AssertCalled(t, "HdelCategory", context.Background(), test.categoryID)
			}
		})
	}
}

func TestCategoryType(t *testing.T) {
	type tests struct {
		Name         string
		typeFound    string
		categories   []domain.CategoryOutput
		mockErr      error
		categoryErr  error
		shouldCallDB bool
	}

	arrTests := []tests{
		{
			Name:      "success",
			typeFound: "траты на еду",
			categories: []domain.CategoryOutput{
				{
					UserID:      1,
					Name:        "food",
					Limit:       10000,
					Type:        "траты на еду",
					Description: "траты на еду",
				},
				{
					UserID:      1,
					Name:        "restaurant",
					Limit:       5000,
					Type:        "траты на еду",
					Description: "ресторан",
				},
			},
			mockErr:      nil,
			categoryErr:  nil,
			shouldCallDB: true,
		},
		{
			Name:         "empty type",
			typeFound:    "",
			categories:   []domain.CategoryOutput{},
			mockErr:      nil,
			categoryErr:  ErrValidateType,
			shouldCallDB: false,
		},
		{
			Name:         "not found",
			typeFound:    "несуществующий тип",
			categories:   []domain.CategoryOutput{},
			mockErr:      postgresql.ErrorNotFound,
			categoryErr:  ErrNoFound,
			shouldCallDB: true,
		},
		{
			Name:         "database error",
			typeFound:    "тип категории",
			categories:   []domain.CategoryOutput{},
			mockErr:      errors.New("database connection error"),
			categoryErr:  ErrDatabase,
			shouldCallDB: true,
		},
	}

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {
			repoMock := new(DbMock)
			r := new(cash.RedisMock)

			if test.shouldCallDB {
				repoMock.On("CategoriesType", mock.Anything, test.typeFound).Return(test.categories, test.mockErr)
			}

			log := logrus.New()

			servic := CreateCategoryServer(repoMock, repoMock, repoMock, repoMock, repoMock, log, r)
			resultCategories, err := servic.CategoryType(context.Background(), test.typeFound)

			if test.categoryErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categoryErr: %v", err)
				}
				assert.Empty(t, resultCategories)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.categories, resultCategories)
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "CategoriesType", mock.Anything, test.typeFound)
			}
			repoMock.AssertExpectations(t)
		})
	}
}
