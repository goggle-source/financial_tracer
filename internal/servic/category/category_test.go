package category

import (
	"errors"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {

	type tests struct {
		Name         string
		userID       uint
		category     domain.CategoryInput
		categoryID   uint
		mockErr      error
		categoryErr  error
		shouldCallDB bool
	}

	arrTests := []tests{
		{
			Name:   "success",
			userID: 1,
			category: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			categoryID:   2,
			mockErr:      nil,
			categoryErr:  nil,
			shouldCallDB: true,
		},
		{
			Name:   "not found",
			userID: 2,
			category: domain.CategoryInput{
				Name:        "default",
				Limit:       10000,
				Description: "test case",
			},
			categoryID:   0,
			mockErr:      postgresql.ErrorNotFound,
			categoryErr:  ErrNoFound,
			shouldCallDB: true,
		},
		{
			Name:   "duplicated",
			userID: 3,
			category: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			categoryID:   0,
			mockErr:      postgresql.ErrorDuplicated,
			categoryErr:  ErrDuplicated,
			shouldCallDB: true,
		},
		{
			Name:   "valdiate",
			userID: 5,
			category: domain.CategoryInput{
				Name:        "a",
				Limit:       100000,
				Description: "asfdasdfgAFG",
			},
			categoryID:   0,
			mockErr:      nil,
			categoryErr:  validator.ValidationErrors{},
			shouldCallDB: false,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("CreateCategoryDatabase", test.userID, test.category).Return(test.categoryID, test.mockErr)

			log := logrus.New()

			servic := CreateCategoryServer(repoMock, log)
			resultID, err := servic.CreateCategory(test.userID, test.category)

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
				repoMock.AssertCalled(t, "CreateCategoryDatabase", test.userID, test.category)
			}
		})
	}
}

func TestReadCategory(t *testing.T) {

	type tests struct {
		Name        string
		category    domain.CategoryOutput
		categoryID  uint
		mockErr     error
		categoryErr error
	}

	arrTests := []tests{
		{
			Name: "success",
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			categoryID:  2,
			mockErr:     nil,
			categoryErr: nil,
		},
		{
			Name:        "not found",
			category:    domain.CategoryOutput{},
			categoryID:  0,
			mockErr:     postgresql.ErrorNotFound,
			categoryErr: ErrNoFound,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("ReadCategoryDatabase", test.categoryID).Return(test.category, test.mockErr)

			log := logrus.New()

			servic := CreateCategoryServer(repoMock, log)
			resultCategory, err := servic.ReadCategory(test.categoryID)

			if test.mockErr != nil || test.categoryErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.category, resultCategory)
			}
			repoMock.AssertCalled(t, "ReadCategoryDatabase", test.categoryID)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestUpdateCategory(t *testing.T) {

	type tests struct {
		Name         string
		newCategory  domain.CategoryInput
		category     domain.CategoryOutput
		categoryID   uint
		mockErr      error
		categoryErr  error
		shouldCallDB bool
	}

	arrTests := []tests{
		{
			Name: "success",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "траты на еду",
			},
			category: domain.CategoryOutput{
				UserID:      2,
				Name:        "food",
				Limit:       15000,
				Description: "на еду и средства",
			},
			categoryID:   2,
			mockErr:      nil,
			categoryErr:  nil,
			shouldCallDB: true,
		},
		{
			Name: "not found",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "ssssss",
			},
			category:     domain.CategoryOutput{},
			categoryID:   0,
			mockErr:      postgresql.ErrorNotFound,
			categoryErr:  ErrNoFound,
			shouldCallDB: true,
		},
		{
			Name: "duplicated",
			newCategory: domain.CategoryInput{
				Name:        "food",
				Limit:       10000,
				Description: "",
			},
			category:     domain.CategoryOutput{},
			categoryID:   1,
			mockErr:      postgresql.ErrorDuplicated,
			categoryErr:  ErrDuplicated,
			shouldCallDB: true,
		},
		{
			Name: "valdiate",
			newCategory: domain.CategoryInput{
				Name:        "",
				Limit:       0,
				Description: "",
			},
			category:     domain.CategoryOutput{},
			categoryID:   0,
			mockErr:      nil,
			categoryErr:  validator.ValidationErrors{},
			shouldCallDB: false,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("UpdateCategoryDatabase", test.categoryID, test.newCategory).Return(test.category, test.mockErr)

			log := logrus.New()

			servic := CreateCategoryServer(repoMock, log)
			_, err := servic.UpdateCategory(test.categoryID, test.newCategory)

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
				repoMock.AssertCalled(t, "UpdateCategoryDatabase", test.categoryID, test.newCategory)
			}
		})
	}
}

func TestDeleteCategoryDatabase(t *testing.T) {
	type tests struct {
		Name        string
		categoryID  uint
		mockErr     error
		categoryErr error
	}

	arrTests := []tests{
		{
			Name:        "success",
			categoryID:  2,
			mockErr:     nil,
			categoryErr: nil,
		},
		{
			Name:        "not found",
			categoryID:  0,
			mockErr:     postgresql.ErrorNotFound,
			categoryErr: ErrNoFound,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("DeleteCategoryDatabase", test.categoryID).Return(test.mockErr)
			log := logrus.New()

			servic := CreateCategoryServer(repoMock, log)
			err := servic.DeleteCategory(test.categoryID)

			if test.mockErr != nil || test.categoryErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, test.categoryErr) {
					t.Fatalf("err != test.categorErr: %v", err)
				}
			} else {
				assert.NoError(t, err)
			}
			repoMock.AssertCalled(t, "DeleteCategoryDatabase", test.categoryID)

			repoMock.AssertExpectations(t)
		})
	}
}
