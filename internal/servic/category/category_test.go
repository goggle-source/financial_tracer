package category

import (
	"errors"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {

	type tests struct {
		Name         string
		userID       uint
		category     domain.CategoryInput
		categoryID   uint
		mokuErr      error
		msgErr       string
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
			mokuErr:      nil,
			msgErr:       "",
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
			mokuErr:      errors.New("error not found"),
			msgErr:       "error create category",
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
			mokuErr:      nil,
			msgErr:       "error validate",
			shouldCallDB: false,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("CreateCategoryDatabase", test.userID, test.category).Return(test.categoryID, test.mokuErr)

			servic := CreateCategoryServer(repoMock)
			resultID, err := servic.CreateCategory(test.userID, test.category)

			if test.mokuErr != nil {
				assert.Error(t, err)
				if test.msgErr != "" {
					assert.Contains(t, err.Error(), test.msgErr)
				}
			}

			assert.Equal(t, test.categoryID, resultID)
			if test.shouldCallDB {
				repoMock.AssertCalled(t, "CreateCategoryDatabase", test.userID, test.category)
			}
		})
	}
}

func TestReadCategory(t *testing.T) {

	type tests struct {
		Name       string
		category   domain.CategoryOutput
		categoryID uint
		mokuErr    error
		msgErr     string
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
			categoryID: 2,
			mokuErr:    nil,
			msgErr:     "",
		},
		{
			Name:       "not found",
			category:   domain.CategoryOutput{},
			categoryID: 0,
			mokuErr:    errors.New("error not found"),
			msgErr:     "error read category",
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("ReadCategoryDatabase", test.categoryID).Return(test.category, test.mokuErr)

			servic := CreateCategoryServer(repoMock)
			resultCategory, err := servic.ReadCategory(test.categoryID)

			if test.mokuErr != nil {
				assert.Error(t, err)
				if test.msgErr != "" {
					assert.Contains(t, err.Error(), test.msgErr)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.category, resultCategory)
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
		mokuErr      error
		msgErr       string
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
			mokuErr:      nil,
			msgErr:       "",
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
			mokuErr:      errors.New("error not found"),
			msgErr:       "error update category",
			shouldCallDB: true,
		},
		{
			Name: "validate",
			newCategory: domain.CategoryInput{
				Name:        "",
				Limit:       0,
				Description: "",
			},
			category:     domain.CategoryOutput{},
			categoryID:   0,
			mokuErr:      nil,
			msgErr:       "error validate",
			shouldCallDB: false,
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("UpdateCategoryDatabase", test.categoryID, test.newCategory).Return(test.category, test.mokuErr)

			servic := CreateCategoryServer(repoMock)
			_, err := servic.UpdateCategory(test.categoryID, test.newCategory)

			if test.mokuErr != nil {
				assert.Error(t, err)
				if test.msgErr != "" {
					assert.Contains(t, err.Error(), test.msgErr)
				}
			}

			if test.shouldCallDB {
				repoMock.AssertCalled(t, "UpdateCategoryDatabase", test.categoryID, test.newCategory)
			}
		})
	}
}

func TestDeleteCategoryDatabase(t *testing.T) {
	type tests struct {
		Name       string
		categoryID uint
		mokuErr    error
		msgErr     string
	}

	arrTests := []tests{
		{
			Name:       "success",
			categoryID: 2,
			mokuErr:    nil,
			msgErr:     "",
		},
		{
			Name:       "not found",
			categoryID: 0,
			mokuErr:    errors.New("error not found"),
			msgErr:     "error delete user",
		},
	}
	repoMock := new(DbMock)

	for _, test := range arrTests {
		t.Run(test.Name, func(t *testing.T) {

			repoMock.On("DeleteCategoryDatabase", test.categoryID).Return(test.mokuErr)

			servic := CreateCategoryServer(repoMock)
			err := servic.DeleteCategory(test.categoryID)

			if test.mokuErr != nil {
				assert.Error(t, err)
				if test.msgErr != "" {
					assert.Contains(t, err.Error(), test.msgErr)
				}
			}
			repoMock.AssertCalled(t, "DeleteCategoryDatabase", test.categoryID)

			repoMock.AssertExpectations(t)
		})
	}
}
