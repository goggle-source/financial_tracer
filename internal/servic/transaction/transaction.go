package transaction

import (
	"context"
	"strconv"

	"github.com/financial_tracer/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CreateTransactionRepository interface {
	CreateTransaction(ctx context.Context, idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
}

type GetTransactionRepository interface {
	GetTransaction(ctx context.Context, TransactionId uint) (domain.TransactionOutput, error)
}

type UpdateTransactionRepository interface {
	UpdateTransaction(ctx context.Context, transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
}

type Redis interface {
	HsetTransaction(ctx context.Context, id uint, transaction domain.TransactionOutput) error
	HgetTransaction(ctx context.Context, id uint) (map[string]string, error)
	HdelTransaction(ctx context.Context, id uint) error
}

type DeleteTransactionRepository interface {
	DeleteTransaction(ctx context.Context, transactionId uint) error
}

type TransactionServer struct {
	d        DeleteTransactionRepository
	c        CreateTransactionRepository
	g        GetTransactionRepository
	u        UpdateTransactionRepository
	log      *logrus.Logger
	validate validator.Validate
	rbd      Redis
}

func CreateTransactionServer(d DeleteTransactionRepository,
	c CreateTransactionRepository,
	g GetTransactionRepository,
	u UpdateTransactionRepository,
	log *logrus.Logger,
	r Redis) *TransactionServer {

	return &TransactionServer{
		d:        d,
		g:        g,
		c:        c,
		u:        u,
		log:      log,
		validate: *validator.New(),
		rbd:      r,
	}
}

func (ts *TransactionServer) CreateTransaction(ctx context.Context, idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	const op = "transaction.CreateTransactionServic"

	log := ts.log.WithFields(logrus.Fields{
		"op":          op,
		"user_id":     idUser,
		"category_id": idCategory,
	})

	log.Info("start create transaction")

	if err := ts.validate.Struct(&tran); err != nil {
		log.WithField("err", err).Error("error validate")

		return 0, err
	}

	id, err := ts.c.CreateTransaction(ctx, idUser, idCategory, tran)
	if err != nil {
		log.Error("error create transaction: ", err)
		return 0, RegisterErrDatabase(err)
	}

	canal := make(chan error)
	go func(canal chan error) {
		transaction := domain.TransactionOutput{
			Name:        tran.Name,
			UserID:      idUser,
			CategoryID:  idCategory,
			Description: tran.Description,
			Count:       tran.Count,
		}
		canal <- ts.rbd.HsetTransaction(ctx, id, transaction)
	}(canal)

	if err := <-canal; err != nil {
		log.Error("error create cash: ", err)
	}

	log.Info("success create transaction")

	return id, nil
}

func (ts *TransactionServer) GetTransaction(ctx context.Context, idTransaction uint) (domain.TransactionOutput, error) {
	const op = "transaction.GetTransactionServer"
	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start read transaction")

	result, err := ts.rbd.HgetTransaction(ctx, idTransaction)
	if err == nil {
		usID, _ := strconv.ParseUint(result["userID"], 10, 64)
		categorID, _ := strconv.ParseUint(result["categoryID"], 10, 64)
		count, _ := strconv.Atoi(result["count"])

		return domain.TransactionOutput{
			Name:        result["name"],
			Description: result["description"],
			UserID:      uint(usID),
			CategoryID:  uint(categorID),
			Count:       count,
		}, nil
	} else {
		log.Info("err info: ", err)
	}

	transaction, err := ts.g.GetTransaction(ctx, idTransaction)
	if err != nil {
		log.Error("error get transaction: ", err)
		return domain.TransactionOutput{}, RegisterErrDatabase(err)
	}

	log.Info("success get transaction")

	return transaction, nil
}

func (ts *TransactionServer) UpdateTransaction(ctx context.Context, idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	const op = "transaction.UpdateTransactionServer"

	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start update transaction")

	if err := ts.validate.Struct(&newTransaction); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return domain.TransactionOutput{}, err
	}

	transaction, err := ts.u.UpdateTransaction(ctx, idTransaction, newTransaction)
	if err != nil {
		log.Error("error update transaction: ", err)
		return domain.TransactionOutput{}, RegisterErrDatabase(err)
	}

	canal := make(chan error)
	go func(canal chan error) {
		canal <- ts.rbd.HsetTransaction(ctx, idTransaction, transaction)
	}(canal)
	if err := <-canal; err != nil {
		log.Error("error update cash: ", err)
	}

	log.Info("success update transaction")

	return transaction, nil
}

func (ts *TransactionServer) DeleteTransaction(ctx context.Context, idTransaction uint) error {
	const op = "transaction.DeleteTransactionServer"

	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start delete transaction")

	err := ts.d.DeleteTransaction(ctx, idTransaction)
	if err != nil {
		log.Error("error delete transaction: ", err)
		return RegisterErrDatabase(err)
	}

	canal := make(chan error)
	go func(canal chan error) {
		canal <- ts.rbd.HdelTransaction(ctx, idTransaction)
	}(canal)

	if err := <-canal; err != nil {
		log.Error("error delete cash: ", err)
	}

	log.Info("success delete transaction")
	return nil
}
