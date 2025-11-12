package categoryHandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/servic/category"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
)

func TestPostCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		userID       uint
		body         domain.CategoryInput
		category     RequestCreateCategory
		categoryID   uint
		status       int
		invalidJSON  bool
		missUserID   bool
		mockErr      error
		shouldCallDB bool
	}{
		{
			name:         "success",
			userID:       1,
			body:         domain.CategoryInput{Name: "food", Limit: 1000, Description: "desc"},
			category:     RequestCreateCategory{Name: "food", Limit: 1000, Description: "desc"},
			categoryID:   10,
			status:       http.StatusOK,
			mockErr:      nil,
			shouldCallDB: true,
		},
		{
			name:         "invalid json",
			userID:       1,
			invalidJSON:  true,
			status:       http.StatusBadRequest,
			mockErr:      nil,
			shouldCallDB: false,
		},
		{
			name:         "no user id",
			missUserID:   true,
			body:         domain.CategoryInput{Name: "food", Limit: 1000, Description: "desc"},
			category:     RequestCreateCategory{Name: "food", Limit: 1000, Description: "desc"},
			categoryID:   0,
			status:       http.StatusInternalServerError,
			mockErr:      nil,
			shouldCallDB: false,
		},
		{
			name:         "error database",
			userID:       12,
			body:         domain.CategoryInput{Name: "credit", Limit: 100000, Description: "max count take a many"},
			category:     RequestCreateCategory{Name: "credit", Limit: 100000, Description: "max count take a many"},
			status:       http.StatusInternalServerError,
			mockErr:      errors.New("error database"),
			shouldCallDB: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if !tc.missUserID {
				c.Set("userID", tc.userID)
			}

			svc := new(categoryServiceMock)
			log := logrus.New()
			ctx := context.Background()

			svc.On("CreateCategory", ctx, tc.userID, tc.body).Return(tc.categoryID, tc.mockErr)
			h := CreateHandlersCategory(svc, svc, svc, svc, svc, log, ctx)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.category)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.PostCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "CreateCategory", ctx, tc.userID, tc.body)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		req          uint
		output       domain.CategoryOutput
		mockErr      error
		status       int
		shouldCallDB bool
	}{
		{
			name:         "success",
			req:          5,
			output:       domain.CategoryOutput{UserID: 1, Name: "food", Limit: 1000, Description: "desc"},
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			req:          5,
			mockErr:      category.ErrNoFound,
			status:       http.StatusNotFound,
			shouldCallDB: true,
		},
		{
			name:         "invalid json",
			status:       http.StatusBadRequest,
			shouldCallDB: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			svc := new(categoryServiceMock)
			log := logrus.New()
			ctx := context.Background()

			if tc.shouldCallDB {
				svc.On("GetCategory", ctx, tc.req).Return(tc.output, tc.mockErr)
			}

			h := CreateHandlersCategory(svc, svc, svc, svc, svc, log, ctx)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}

			b, _ := json.Marshal(tc.req)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

			req.Header.Set("content-type", "application/json")
			c.Request = &req
			var strID string

			if tc.req == 0 {
				strID = "adsa"
			} else {
				strID = strconv.Itoa(int(tc.req))
			}
			c.Params = gin.Params{
				{
					Key:   "id",
					Value: strID,
				},
			}

			h.GetCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "GetCategory", ctx, tc.req)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		req          RequestUpdateCategory
		output       domain.CategoryOutput
		mockErr      error
		status       int
		invalidJSON  bool
		shouldCallDB bool
	}{
		{
			name:         "success",
			req:          RequestUpdateCategory{Name: "new", Limit: 200, Description: "d", CategoryId: 7},
			output:       domain.CategoryOutput{UserID: 1, Name: "new", Limit: 200, Description: "d"},
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			req:          RequestUpdateCategory{Name: "new", Limit: 200, Description: "d", CategoryId: 77},
			mockErr:      category.ErrNoFound,
			status:       http.StatusNotFound,
			shouldCallDB: true,
		},
		{
			name:         "invalid json",
			invalidJSON:  true,
			status:       http.StatusBadRequest,
			shouldCallDB: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			svc := new(categoryServiceMock)
			log := logrus.New()
			ctx := context.Background()
			input := domain.CategoryInput{Name: tc.req.Name, Limit: tc.req.Limit, Description: tc.req.Description}

			if tc.shouldCallDB {
				svc.On("UpdateCategory", ctx, tc.req.CategoryId, input).Return(tc.output, tc.mockErr)
			}
			h := CreateHandlersCategory(svc, svc, svc, svc, svc, log, ctx)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.req)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.UpdateCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "UpdateCategory", ctx, tc.req.CategoryId, input)
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		req          uint
		mockErr      error
		status       int
		shouldCallDB bool
	}{
		{
			name:         "success",
			req:          5,
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			req:          55,
			mockErr:      category.ErrNoFound,
			status:       http.StatusNotFound,
			shouldCallDB: true,
		},
		{
			name:         "invalid id",
			req:          0,
			status:       http.StatusBadRequest,
			shouldCallDB: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			svc := new(categoryServiceMock)
			log := logrus.New()
			ctx := context.Background()

			if tc.shouldCallDB {
				svc.On("DeleteCategory", ctx, tc.req).Return(tc.mockErr)
			}

			h := CreateHandlersCategory(svc, svc, svc, svc, svc, log, ctx)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}

			b, _ := json.Marshal(tc.req)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

			req.Header.Set("content-type", "application/json")
			c.Request = &req
			var strID string
			if tc.req == 0 {
				strID = "asdad"
			} else {
				strID = strconv.Itoa(int(tc.req))
			}

			c.Params = gin.Params{
				{
					Key:   "id",
					Value: strID,
				},
			}

			h.DeleteCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "DeleteCategory", ctx, tc.req)
			}
		})
	}
}

func TestCategoryType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		param        string
		output       []domain.CategoryOutput
		mockErr      error
		status       int
		shouldCallDB bool
	}{
		{
			name:  "success",
			param: "income",
			output: []domain.CategoryOutput{
				{UserID: 1, Name: "Salary", Limit: 50000, Type: "income", Description: "Monthly salary"},
				{UserID: 1, Name: "Freelance", Limit: 10000, Type: "income", Description: "Freelance work"},
			},
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "empty param",
			param:        "",
			status:       http.StatusBadRequest,
			shouldCallDB: false,
		},
		{
			name:         "invalid type",
			param:        "invalid",
			mockErr:      category.ErrValidateType,
			status:       http.StatusBadRequest,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			param:        "expense",
			mockErr:      category.ErrNoFound,
			status:       http.StatusNotFound,
			shouldCallDB: true,
		},
		{
			name:         "database error",
			param:        "income",
			mockErr:      category.ErrDatabase,
			status:       http.StatusInternalServerError,
			shouldCallDB: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			svc := new(categoryServiceMock)
			log := logrus.New()
			ctx := context.Background()

			if tc.shouldCallDB {
				svc.On("CategoryType", ctx, tc.param).Return(tc.output, tc.mockErr)
			}

			h := CreateHandlersCategory(svc, svc, svc, svc, svc, log, ctx)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			if tc.param != "" {
				c.Params = gin.Params{
					{
						Key:   "type",
						Value: tc.param,
					},
				}
			}

			h.CategoryType(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "CategoryType", ctx, tc.param)
			}
		})
	}
}
