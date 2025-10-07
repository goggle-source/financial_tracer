package userHandlers

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
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
)

func TestRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         any
		user         domain.RegisterUser
		userID       uint
		userName     string
		status       int
		invalidJSON  bool
		mockErr      error
		shouldCallDB bool
	}{
		{
			name: "success",
			body: UserRegistration{
				Name:        "Alice Doe",
				UserRequest: UserRequest{Email: "alice@example.com", Password: "secret"},
			},
			user:         domain.RegisterUser{Name: "Alice Doe", Email: "alice@example.com", Password: "secret"},
			userID:       1,
			userName:     "Alice Doe",
			status:       http.StatusOK,
			mockErr:      nil,
			shouldCallDB: true,
		},
		{
			name: "duplicated",
			body: UserRegistration{
				Name:        "Bob",
				UserRequest: UserRequest{Email: "bob@example.com", Password: "qwerty"},
			},
			user:         domain.RegisterUser{Name: "Bob", Email: "bob@example.com", Password: "qwerty"},
			userID:       0,
			userName:     "",
			status:       http.StatusBadRequest,
			mockErr:      postgresql.ErrorDuplicated,
			shouldCallDB: true,
		},
		{
			name:         "invalid json",
			body:         nil,
			status:       http.StatusBadRequest,
			invalidJSON:  true,
			shouldCallDB: false,
		},
		{
			name: "internal error",
			body: UserRegistration{
				Name:        "Alice Doe",
				UserRequest: UserRequest{Email: "alice@example.com", Password: "secret"},
			},
			user:         domain.RegisterUser{Name: "Alice Doe", Email: "alice@example.com", Password: "secret"},
			userID:       0,
			userName:     "",
			status:       http.StatusInternalServerError,
			mockErr:      errors.New("error database"),
			shouldCallDB: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			svc := new(userServiceMock)
			log := logrus.New()

			svc.On("ServerRegistrationUser", tc.user).Return(tc.userID, tc.userName, tc.mockErr)

			h := CreateHandlersUser("secret", svc, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.Registration(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "ServerRegistrationUser", tc.user)
			}
		})
	}
}

func TestAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         any
		mockRet      func(*userServiceMock)
		user         domain.AuthenticationUser
		userID       uint
		userName     string
		status       int
		invalidJSON  bool
		mockErr      error
		shouldCallDB bool
	}{
		{
			name:         "success",
			body:         UserRequest{Email: "alice@example.com", Password: "secret"},
			user:         domain.AuthenticationUser{Email: "alice@example.com", Password: "secret"},
			userID:       1,
			userName:     "alice",
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			body:         UserRequest{Email: "missing@example.com", Password: "x"},
			user:         domain.AuthenticationUser{Email: "missing@example.com", Password: "x"},
			userID:       0,
			userName:     "",
			status:       http.StatusBadRequest,
			mockErr:      postgresql.ErrorNotFound,
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

			svc := new(userServiceMock)
			log := logrus.New()

			svc.On("ServerAuthenticationUser", tc.user).Return(tc.userID, tc.userName, tc.mockErr)

			h := CreateHandlersUser("secret", svc, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.Authentication(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "ServerAuthenticationUser", tc.user)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         any
		user         domain.DeleteUser
		mockErr      error
		status       int
		invalidJSON  bool
		shouldCallDB bool
	}{
		{
			name:         "success",
			body:         UserRequest{Email: "alice@example.com", Password: "secret"},
			user:         domain.DeleteUser{Email: "alice@example.com", Password: "secret"},
			mockErr:      nil,
			status:       http.StatusOK,
			shouldCallDB: true,
		},
		{
			name:         "not found",
			body:         UserRequest{Email: "missing@example.com", Password: "x"},
			user:         domain.DeleteUser{Email: "missing@example.com", Password: "x"},
			mockErr:      postgresql.ErrorNotFound,
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

			svc := new(userServiceMock)
			log := logrus.New()

			svc.On("ServerDeleteUser", tc.user).Return(tc.mockErr)

			h := CreateHandlersUser("secret", svc, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				b, _ := json.Marshal(tc.body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			h.DeleteUser(c)
			assert.Equal(t, tc.status, w.Code)
			if tc.shouldCallDB {
				svc.AssertCalled(t, "ServerDeleteUser", tc.user)
			}
		})
	}
}

func TestGetAccessToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	log := logrus.New()
	h := CreateHandlersUser("secret", nil, log)

	// invalid json
	req := http.Request{Header: make(http.Header), URL: &url.URL{}}
	req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
	req.Header.Set("content-type", "application/json")
	c.Request = &req

	h.GetAccessToken(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
