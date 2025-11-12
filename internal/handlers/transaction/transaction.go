package transactionHandlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/handlers/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CreateTransactionServic interface {
	CreateTransaction(ctx context.Context, idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
}

type GetTransactionServic interface {
	GetTransaction(ctx context.Context, idTransaction uint) (domain.TransactionOutput, error)
}

type UpdateTransactionServic interface {
	UpdateTransaction(ctx context.Context, idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
}

type DeleteTransactionServic interface {
	DeleteTransaction(ctx context.Context, idTransaction uint) error
}

type TransactionHandlers struct {
	c   CreateTransactionServic
	g   GetTransactionServic
	u   UpdateTransactionServic
	d   DeleteTransactionServic
	log *logrus.Logger
	ctx context.Context
}

func CreateTransactionHandlers(c CreateTransactionServic,
	g GetTransactionServic,
	u UpdateTransactionServic,
	d DeleteTransactionServic,
	log *logrus.Logger,
	ctx context.Context) *TransactionHandlers {
	return &TransactionHandlers{
		c:   c,
		d:   d,
		g:   g,
		u:   u,
		log: log,
		ctx: ctx,
	}
}

// CreateTransaction godoc
//
//	@Summary		Создание транзакции
//	@Description	Создание 1 транзакции для пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestCreateTransaction	true	"данные для создание пользователя"
//	@Success		200	{object}	api.SuccessResponse			"Транзакция создана успешно"
//
//	@Failure		401	{object}	api.ErrorResponse			"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse			"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse			"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse			"Некорректные данные"
//	@Router			/transaction/ [post]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) PostTransaction(c *gin.Context) {
	const op = "handlers.PostTransaction"

	log := th.log.WithField("op", op)

	log.Info("start create transaction")

	var transaction RequestCreateTransaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.WithField("err", err).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}
	idUser, ok := c.Get("userID")
	if !ok {
		log.Error("error get userID")
		api.ResponseError(c, http.StatusInternalServerError, "error server")
		return
	}

	newTransaction := domain.TransactionInput{
		Name:        transaction.Name,
		Count:       transaction.Count,
		Description: transaction.Description,
	}

	id, err := th.c.CreateTransaction(th.ctx, idUser.(uint), transaction.IdCategory, newTransaction)
	if err != nil {
		log.WithField("err", err).Error("error create transaction")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success create transaction")

	api.ResponseOK(c, id)
}

// GetTransaction godoc
//
//	@Summary		Получение транзакции
//	@Description	Получение 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Produce		json
//	@Param			id	path		int					true	"id транзакции"
//	@Success		200	{object}	api.SuccessResponse	"good"
//
//	@Failure		401	{object}	api.ErrorResponse	"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse	"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse	"Транзакция не найдена"
//	@Router			/transaction/{id} [get]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) GetTransaction(c *gin.Context) {
	const op = "handlers.GetTransaction"

	log := th.log.WithField("op", op)

	log.Info("start get transaction")

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.WithField("err", err).Error("invalid convert string in int", err)
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	transaction, err := th.g.GetTransaction(c.Request.Context(), uint(id))
	if err != nil {
		log.WithField("err", err).Error("error get transaction")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success get transaction")

	api.ResponseOK(c, transaction)
}

// UpdateTransaction godoc
//
//	@Summary		Обновление транзакции
//	@Description	Обновление 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			req	body		RequestUpdateTransaction	true	"Данные для обновление пользователя"
//	@Success		200	{object}	api.SuccessResponse			"Транзакция обновлена"
//
//	@Failure		401	{object}	api.ErrorResponse			"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse			"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse			"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse			"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse			"Транзакция не найдена"
//	@Router			/transaction/ [put]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) UpdateTransaction(c *gin.Context) {
	const op = "handlers.UpdateTransaction"

	log := th.log.WithField("op", op)

	log.Info("start update transaction")

	var transaction RequestUpdateTransaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.WithField("err", err).Error("error valid JSON")
		api.ResponseError(c, http.StatusBadRequest, "error valid JSON")
		return
	}
	tran := domain.TransactionInput{
		Name:        transaction.Name,
		Count:       transaction.Count,
		Description: transaction.Description,
	}

	newTransaction, err := th.u.UpdateTransaction(c.Request.Context(), transaction.IdTransaction, tran)
	if err != nil {
		log.WithField("err", err).Error("error update transaction")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success update transaction")

	api.ResponseOK(c, newTransaction)
}

// DeleteTransaction godoc
//
//	@Summary		Удаление транзакции
//	@Description	Удаление 1 транзакции для 1 пользователя
//	@Tags			transaction
//	@Produce		json
//	@Param			id	path		int					true	"id транзакции"
//	@Success		200	{object}	api.SuccessResponse	"Транзакция удалена"
//
//	@Failure		401	{object}	api.ErrorResponse	"Ошибка авторизации"
//
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные входные данные"
//	@Failure		500	{object}	api.ErrorResponse	"Ошибка сервера"
//	@Failure		400	{object}	api.ErrorResponse	"Некорректные данные"
//	@Failure		404	{object}	api.ErrorResponse	"Транзакция не найдена"
//	@Router			/transaction/{id} [delete]
//
//	@Security		jwtAuth
func (th *TransactionHandlers) DeleteTransaction(c *gin.Context) {
	const op = "handlers.DeleteTransaction"

	log := th.log.WithField("op", op)

	log.Info("start delete transaction")

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.WithField("err", err).Error("invalid convert string in int", err)
		api.ResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	err = th.d.DeleteTransaction(c.Request.Context(), uint(id))

	if err != nil {
		log.WithField("err", err).Error("error delete transaction")
		api.RegistrationError(c, err)
		return
	}

	log.Info("success delete transaction")

	api.ResponseOK(c, "transaction delete")
}
