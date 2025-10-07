package transactionHandlers

import (
	"net/http"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/handlers/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServicTransactier interface {
	CreateTransactionServic(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
	ReadTransactionServer(idTransaction uint) (domain.TransactionOutput, error)
	UpdateTransactionServer(idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
	DeleteTransactionServer(idTransaction uint) error
}

type TransactionHandlers struct {
	transaction ServicTransactier
	log         *logrus.Logger
}

func CreateTransactionHandlers(s ServicTransactier, log *logrus.Logger) *TransactionHandlers {
	return &TransactionHandlers{
		transaction: s,
		log:         log,
	}
}

// CreateTransaction godoc
//
//	@Summary		Создание транзакции
//	@Description	Создание 1 транзакции для пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestCreateTransaction				true	"данные для создание пользователя"
//	@Success		200	{object}	api.SuccessResponse[uint]				"Транзакция создана успешно"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Router			/transaction/create [post]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) PostTransaction(c *gin.Context) {
	const op = "handlers.PostTransaction"
	var transaction RequestCreateTransaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}
	idUser, ok := c.Get("userID")
	if !ok {
		th.log.WithField("op", op).Error("error get userID")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	newTransaction := domain.TransactionInput{
		Name:        transaction.Name,
		Count:       transaction.Count,
		Description: transaction.Description,
	}

	id, err := th.transaction.CreateTransactionServic(idUser.(uint), transaction.IdCategory, newTransaction)
	if err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error create transaction")
		api.RegistrationError(c, op, err)
		return
	}
	api.ResponseOK(c, id)
}

// GetTransaction godoc
//
//	@Summary		Получение транзакции
//	@Description	Получение 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestIdTransaction					true	"id транзакции"
//	@Success		200	{object}	api.SuccessResponse[domain.Transaction]	"good"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Транзакция не найдена"
//	@Router			/transaction/get [get]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) GetTransaction(c *gin.Context) {
	const op = "handlers.GetTransaction"
	var id RequestIdTransaction

	if err := c.ShouldBindJSON(&id); err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	transaction, err := th.transaction.ReadTransactionServer(id.IdTransaction)
	if err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error get transaction")
		api.RegistrationError(c, op, err)
		return
	}

	api.ResponseOK(c, transaction)
}

// UpdateTransaction godoc
//
//	@Summary		Обновление транзакции
//	@Description	Обновление 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestUpdateTransaction				true	"Данные для обновление пользователя"
//	@Success		200	{object}	api.SuccessResponse[domain.Transaction]	"Транзакция обновлена"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Транзакция не найдена"
//	@Router			/transaction/update [put]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) UpdateTransaction(c *gin.Context) {
	const op = "handlers.UpdateTransaction"

	var transaction RequestUpdateTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}
	tran := domain.TransactionInput{
		Name:        transaction.Name,
		Count:       transaction.Count,
		Description: transaction.Description,
	}

	newTransaction, err := th.transaction.UpdateTransactionServer(transaction.IdTransaction, tran)
	if err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error update transaction")
		api.RegistrationError(c, op, err)
		return
	}

	api.ResponseOK(c, newTransaction)
}

// DeleteTransaction godoc
//
//	@Summary		Удаление транзакции
//	@Description	Удаление 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestIdTransaction					true	"id транзакции"
//	@Success		200	{object}	api.SuccessResponse[string]				"Транзакция удалена"
//
//	@Failure		401	{object}	api.ErrorResponse[string]				"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse[[]map[string]string]	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse[string]				"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse[string]				"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse[string]				"Транзакция не найдена"
//	@Router			/transaction/delete [delete]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) DeleteTransaction(c *gin.Context) {
	const op = "handlers.DeleteTransaction"
	var id RequestIdTransaction

	if err := c.ShouldBindJSON(&id); err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}

	err := th.transaction.DeleteTransactionServer(id.IdTransaction)
	if err != nil {
		th.log.WithFields(logrus.Fields{
			"op":  op,
			"err": err,
		}).Error("error delete transaction")
		api.RegistrationError(c, op, err)
		return
	}

	api.ResponseOK(c, "transaction delete")
}
