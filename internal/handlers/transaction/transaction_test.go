package transactionHandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/servic/transaction"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
)

func TestCreateTransactionServic(t *testing.T) {
	type test struct {
		name          string
		tran          RequestCreateTransaction
		idUser        uint
		idCategory    uint
		idTransaction uint
		mockErr       error
		status        int
		shouldCallDB  bool
		invalidJSON   bool
	}

	arrTest := []test{
		{
			name: "success",
			tran: RequestCreateTransaction{
				IdCategory:  1,
				Name:        "продукты",
				Count:       1000,
				Description: "покупка продуктов",
			},
			idUser:        1,
			idCategory:    1,
			idTransaction: 1,
			mockErr:       nil,
			status:        http.StatusOK,
			shouldCallDB:  true,
		},
		{
			name: "error database",
			tran: RequestCreateTransaction{
				IdCategory:  1,
				Name:        "продукты",
				Count:       1000,
				Description: "покупка продуктов",
			},
			idUser:        1,
			idCategory:    1,
			idTransaction: 0,
			mockErr:       errors.New("error database"),
			status:        http.StatusInternalServerError,
			shouldCallDB:  true,
		},
		{
			name:          "error validate",
			tran:          RequestCreateTransaction{},
			idUser:        1,
			idCategory:    0,
			idTransaction: 0,
			mockErr:       nil,
			status:        http.StatusBadRequest,
			shouldCallDB:  false,
			invalidJSON:   true,
		},
		{
			name: "error not found",
			tran: RequestCreateTransaction{
				IdCategory:  3,
				Name:        "продукты",
				Count:       1000,
				Description: "покупка продуктов",
			},
			idUser:        3,
			idCategory:    3,
			idTransaction: 0,
			mockErr:       transaction.ErrNoFound,
			status:        http.StatusNotFound,
			shouldCallDB:  true,
		},
		{
			name: "error limit",
			tran: RequestCreateTransaction{
				IdCategory:  2,
				Name:        "продукты",
				Count:       1000,
				Description: "покупка продуктов",
			},
			idUser:        1,
			idCategory:    2,
			idTransaction: 0,
			mockErr:       transaction.ErrLimit,
			status:        http.StatusBadRequest,
			shouldCallDB:  true,
		},
	}

	for _, ts := range arrTest {

		t.Run(ts.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("userID", ts.idUser)

			repoMock := new(tranasctionServicMock)
			log := logrus.New()
			tranInput := domain.TransactionInput{
				Name:        ts.tran.Name,
				Count:       ts.tran.Count,
				Description: ts.tran.Description,
			}

			repoMock.On("CreateTransactionServic", ts.idUser, ts.idCategory, tranInput).Return(ts.idTransaction, ts.mockErr)

			handler := CreateTransactionHandlers(repoMock, repoMock, repoMock, repoMock, log)

			req := http.Request{
				Header: make(http.Header),
				URL:    &url.URL{},
			}

			js, _ := json.Marshal(ts.tran)
			if ts.invalidJSON {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				req.Body = ioutil.NopCloser(bytes.NewBuffer(js))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			handler.PostTransaction(c)

			assert.Equal(t, ts.status, w.Code)

			if ts.shouldCallDB {
				repoMock.AssertCalled(t, "CreateTransactionServic", ts.idUser, ts.idCategory, tranInput)
			}
		})
	}
}

func TestGetTransaction(t *testing.T) {
	type test struct {
		name    string
		req     uint
		output  domain.TransactionOutput
		mockErr error
		status  int
		invalid bool
	}

	cases := []test{
		{
			name:    "success",
			req:     10,
			output:  domain.TransactionOutput{UserID: 1, CategoryID: 2, Name: "food", Count: 100, Description: "desc"},
			mockErr: nil,
			status:  http.StatusOK,
		},
		{
			name:    "not found",
			req:     99,
			output:  domain.TransactionOutput{},
			mockErr: transaction.ErrNoFound,
			status:  http.StatusNotFound,
		},
		{
			name:    "invalid id",
			req:     0,
			invalid: true,
			status:  http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			repoMock := new(tranasctionServicMock)
			log := logrus.New()

			if !tc.invalid {
				repoMock.On("ReadTransactionServer", tc.req).Return(tc.output, tc.mockErr)
			}

			handler := CreateTransactionHandlers(repoMock, repoMock, repoMock, repoMock, log)

			req := http.Request{Header: make(http.Header), URL: &url.URL{}}

			body, _ := json.Marshal(tc.req)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			req.Header.Set("content-type", "application/json")
			c.Request = &req
			strID := strconv.Itoa(int(tc.req))

			if tc.req == 0 {
				strID = "asdas"
			}

			c.Params = gin.Params{
				{
					Key:   "id",
					Value: strID,
				},
			}

			handler.GetTransaction(c)
			assert.Equal(t, tc.status, w.Code)

			if !tc.invalid {
				repoMock.AssertCalled(t, "ReadTransactionServer", tc.req)
			}
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	type test struct {
		name    string
		req     RequestUpdateTransaction
		output  domain.TransactionOutput
		mockErr error
		status  int
		invalid bool
	}

	cases := []test{
		{
			name:    "success",
			req:     RequestUpdateTransaction{IdTransaction: 7, Name: "taxi", Count: 200, Description: "city"},
			output:  domain.TransactionOutput{UserID: 1, CategoryID: 2, Name: "taxi", Count: 200, Description: "city"},
			mockErr: nil,
			status:  http.StatusOK,
		},
		{
			name:    "not found",
			req:     RequestUpdateTransaction{IdTransaction: 77, Name: "taxi", Count: 200, Description: "city"},
			output:  domain.TransactionOutput{},
			mockErr: transaction.ErrNoFound,
			status:  http.StatusNotFound,
		},
		{
			name:    "invalid json",
			invalid: true,
			status:  http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			repoMock := new(tranasctionServicMock)
			log := logrus.New()

			if !tc.invalid {
				input := domain.TransactionInput{Name: tc.req.Name, Count: tc.req.Count, Description: tc.req.Description}
				repoMock.On("UpdateTransactionServer", tc.req.IdTransaction, input).Return(tc.output, tc.mockErr)
			}

			handler := CreateTransactionHandlers(repoMock, repoMock, repoMock, repoMock, log)
			req := http.Request{Header: make(http.Header), URL: &url.URL{}}
			if tc.invalid {
				req.Body = ioutil.NopCloser(bytes.NewBufferString("{"))
			} else {
				body, _ := json.Marshal(tc.req)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
			req.Header.Set("content-type", "application/json")
			c.Request = &req

			handler.UpdateTransaction(c)
			assert.Equal(t, tc.status, w.Code)

			if !tc.invalid {
				input := domain.TransactionInput{Name: tc.req.Name, Count: tc.req.Count, Description: tc.req.Description}
				repoMock.AssertCalled(t, "UpdateTransactionServer", tc.req.IdTransaction, input)
			}
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	type test struct {
		name    string
		req     uint
		mockErr error
		status  int
		invalid bool
	}

	cases := []test{
		{
			name:    "success",
			req:     5,
			mockErr: nil,
			status:  http.StatusOK,
			invalid: false,
		},
		{
			name:    "not found",
			req:     55,
			mockErr: transaction.ErrNoFound,
			status:  http.StatusNotFound,
			invalid: false,
		},
		{
			name:    "invalid id",
			req:     0,
			status:  http.StatusBadRequest,
			invalid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			repoMock := new(tranasctionServicMock)
			log := logrus.New()

			if !tc.invalid {
				repoMock.On("DeleteTransactionServer", tc.req).Return(tc.mockErr)
			}

			handler := CreateTransactionHandlers(repoMock, repoMock, repoMock, repoMock, log)
			req := http.Request{Header: make(http.Header), URL: &url.URL{}}

			req.Header.Set("content-type", "application/json")
			c.Request = &req
			strID := strconv.Itoa(int(tc.req))

			if tc.req == 0 {
				strID = "asdas"
			}

			c.Params = gin.Params{
				{
					Key:   "id",
					Value: strID,
				},
			}

			handler.DeleteTransaction(c)
			assert.Equal(t, tc.status, w.Code)

			if !tc.invalid {
				repoMock.AssertCalled(t, "DeleteTransactionServer", tc.req)
			}
		})
	}
}
