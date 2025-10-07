package categoryHandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/financial_tracer/internal/domain"
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

			svc.On("CreateCategory", tc.userID, tc.body).Return(tc.categoryID, tc.mockErr)
			h := CreateHandlersCategory(svc, log)

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
				svc.AssertCalled(t, "CreateCategory", tc.userID, tc.body)
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		req          IDCategory
		output       domain.CategoryOutput
		mockErr      error
		status       int
		invalidJSON  bool
		shouldCallDB bool
	}{
		{
			name:         "success",
			req:          IDCategory{CategoryId: 5},
			output:       domain.CategoryOutput{UserID: 1, Name: "food", Limit: 1000, Description: "desc"},
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			req:          IDCategory{CategoryId: 55},
			mockErr:      errors.New("error not found"),
			status:       http.StatusBadRequest,
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
			svc.On("ReadCategory", tc.req.CategoryId).Return(tc.output, tc.mockErr)
			h := CreateHandlersCategory(svc, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.req)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.GetCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "ReadCategory", tc.req.CategoryId)
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
			mockErr:      errors.New("error not found"),
			status:       http.StatusBadRequest,
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
			input := domain.CategoryInput{Name: tc.req.Name, Limit: tc.req.Limit, Description: tc.req.Description}

			svc.On("UpdateCategory", tc.req.CategoryId, input).Return(tc.output, tc.mockErr)
			h := CreateHandlersCategory(svc, log)

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
				svc.AssertCalled(t, "UpdateCategory", tc.req.CategoryId, input)
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		req          IDCategory
		mockErr      error
		status       int
		invalidJSON  bool
		shouldCallDB bool
	}{
		{
			name:         "success",
			req:          IDCategory{CategoryId: 5},
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			req:          IDCategory{CategoryId: 55},
			mockErr:      errors.New("error not found"),
			status:       http.StatusBadRequest,
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
			svc.On("DeleteCategory", tc.req.CategoryId).Return(tc.mockErr)
			h := CreateHandlersCategory(svc, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.req)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.DeleteCategory(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "DeleteCategory", tc.req.CategoryId)
			}
		})
	}
}
