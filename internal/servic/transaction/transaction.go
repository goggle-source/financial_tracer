package transaction

import (
	"errors"
	"fmt"

	"github.com/financial_tracer/internal/domain"
	"github.com/financial_tracer/internal/infastructure/db/postgresql"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type DatabaseTransactionRepository interface {
	CreateTransaction(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error)
	GetTransaction(TransactionId uint) (domain.TransactionOutput, error)
	UpdateTransaction(transactionId uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error)
	DeleteTransaction(transactionId uint) error
}

type TransactionServer struct {
	d   DatabaseTransactionRepository
	log *logrus.Logger
}

func CreateTransactionServer(d DatabaseTransactionRepository, log *logrus.Logger) *TransactionServer {
	return &TransactionServer{
		d:   d,
		log: log,
	}
}

func (ts *TransactionServer) CreateTransactionServic(idUser uint, idCategory uint, tran domain.TransactionInput) (uint, error) {
	const op = "transaction.CreateTransactionServic"

	log := ts.log.WithFields(logrus.Fields{
		"op":      op,
		"user_id": idUser,
		"category_id": idCategory,
	})

	log.Info("start create transaction")

	if err := validator.New().Struct(&tran); err != nil {
		log.WithField("err", err).Error("error validate")

		return 0, fmt.Errorf("%s error validate: %w", op, err)
	}

	id, err := ts.d.CreateTransaction(idUser, idCategory, tran)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("error create transaction")
			return 0, fmt.Errorf("%s error create transaction: %w", op, ErrNoFound)
		}
		if errors.Is(err, postgresql.ErrorLimit) {
			log.WithField("err", err).Error("error create transaction")
			return 0, fmt.Errorf("%s error create transaction: %w", op, ErrLimit)
		}
		log.WithField("err", err).Error("error create transaction")

		return 0, fmt.Errorf("%s error create transaction: %w", op, ErrDatabase)
	}
	log.Info("success create transaction")

	return id, nil
}

func (ts *TransactionServer) ReadTransactionServer(idTransaction uint) (domain.TransactionOutput, error) {
	const op = "transaction.ReadTransactionServer"
	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start read transaction")

	transaction, err := ts.d.GetTransaction(idTransaction)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("error get transaction")

			return domain.TransactionOutput{}, fmt.Errorf("%s error get transaction: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("error get transaction")

		return domain.TransactionOutput{}, fmt.Errorf("%s error get transaction: %w", op, ErrDatabase)
	}

	log.Info("success read transaction")

	return transaction, nil
}

func (ts *TransactionServer) UpdateTransactionServer(idTransaction uint, newTransaction domain.TransactionInput) (domain.TransactionOutput, error) {
	const op = "transaction.UpdateTransactionServer"

	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start update transaction")

	if err := validator.New().Struct(&newTransaction); err != nil {
		log.WithField("err", err).Error("invalid validate")

		return domain.TransactionOutput{}, fmt.Errorf("%s invalid validate: %w", op, err)
	}

	transaction, err := ts.d.UpdateTransaction(idTransaction, newTransaction)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("transaction is not found")

			return domain.TransactionOutput{}, fmt.Errorf("%s transaction is not found: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("field update transaction")

		return domain.TransactionOutput{}, fmt.Errorf("%s field update transaction: %w", op, ErrDatabase)
	}
	log.Info("success update transaction")

	return transaction, nil
}

func (ts *TransactionServer) DeleteTransactionServer(idTransaction uint) error {
	const op = "transaction.DeleteTransactionServer"

	log := ts.log.WithFields(logrus.Fields{
		"op":             op,
		"transaction_id": idTransaction,
	})

	log.Info("start delete transaction")

	err := ts.d.DeleteTransaction(idTransaction)
	if err != nil {
		if errors.Is(err, postgresql.ErrorNotFound) {
			log.WithField("err", err).Error("error delete transaction")

			return fmt.Errorf("%s error delete transaction: %w", op, ErrNoFound)
		}
		log.WithField("err", err).Error("error delete transaction")

		return fmt.Errorf("%s error delete transaction: %w", op, ErrDatabase)
	}

	log.Info("success delete transaction")
	return nil
}
