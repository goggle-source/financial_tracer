package transaction

import (
	"github.com/financial_tracer/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CreateTransactionRepository interface {
	CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
}

type GetTransactionRepository interface {
	GetTransaction(TransactionId uint) (domain.TransactionOutput, error)
}

type UpdateTransactionRepository interface {
	UpdateTransaction(transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
}

type DeleteTransactionRepository interface {
	DeleteTransaction(transactionId uint) error
}

type TransactionServer struct {
	d        DeleteTransactionRepository
	c        CreateTransactionRepository
	g        GetTransactionRepository
	u        UpdateTransactionRepository
	log      *logrus.Logger
	validate validator.Validate
}

func CreateTransactionServer(d DeleteTransactionRepository,
	c CreateTransactionRepository,
	g GetTransactionRepository,
	u UpdateTransactionRepository,
	log *logrus.Logger) *TransactionServer {

	return &TransactionServer{
		d:        d,
		g:        g,
		c:        c,
		u:        u,
		log:      log,
		validate: *validator.New(),
	}
}

func (ts *TransactionServer) CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
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

	id, err := ts.c.CreateTransaction(idUser, idCategory, tran)
	if err != nil {
		log.Error("error create transaction: ", err)
		return 0, RegisterErrDatabase(err)
	}
	log.Info("success create transaction")

	return id, nil
}

func (ts *TransactionServer) GetTransaction(idTransaction uint) (domain.TransactionOutput, error) {
	const op = "transaction.GetTransactionServer"
	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start read transaction")

	transaction, err := ts.g.GetTransaction(idTransaction)
	if err != nil {
		log.Error("error get transaction: ", err)
		return domain.TransactionOutput{}, RegisterErrDatabase(err)
	}

	log.Info("success get transaction")

	return transaction, nil
}

func (ts *TransactionServer) UpdateTransaction(idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
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

	transaction, err := ts.u.UpdateTransaction(idTransaction, newTransaction)
	if err != nil {
		log.Error("error update transaction: ", err)
		return domain.TransactionOutput{}, RegisterErrDatabase(err)
	}
	log.Info("success update transaction")

	return transaction, nil
}

func (ts *TransactionServer) DeleteTransaction(idTransaction uint) error {
	const op = "transaction.DeleteTransactionServer"

	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start delete transaction")

	err := ts.d.DeleteTransaction(idTransaction)
	if err != nil {
		log.Error("error delete transaction: ", err)
		return RegisterErrDatabase(err)
	}

	log.Info("success delete transaction")
	return nil
}
