package usecases

import (
	"fmt"
	"rinha-backend-2024-q1/application/repositories"
	"rinha-backend-2024-q1/model"
	"time"
)

type CreateTransactionUseCase struct {
	repository repositories.ClientRepository
}

type CreateTransactionResponse struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

func NewCreateTransactionUseCase(repository repositories.ClientRepository) CreateTransactionUseCase {
	return CreateTransactionUseCase{
		repository: repository,
	}
}

func (u CreateTransactionUseCase) Execute(
	clientId int,
	value int64,
	transactionType string,
	description string,
	clientLimit int64) (CreateTransactionResponse, error) {
	statement, err := u.repository.GetStatement(clientId)
	if err != nil {
		return CreateTransactionResponse{}, nil
	}
	total := statement.Summary.Total
	if transactionType == "d" && futureValueLessThanLimit(value, total, clientLimit) {
		return CreateTransactionResponse{}, fmt.Errorf("LIMIT_NOT_ALLOWED")
	}
	newTransaction := model.Transaction{
		Value:           value,
		TransactionType: transactionType,
		Description:     description,
		Timestamp:       time.Now().String(),
	}
	transactions := append([]model.Transaction{newTransaction}, statement.Transactions...)
	if len(transactions) > 10 {
		transactions = transactions[:len(transactions)-1]
	}
	statement.Transactions = transactions
	if transactionType == "d" {
		statement.Summary.Total = total - value
	} else {
		statement.Summary.Total = total + value
	}
	err = u.repository.SaveStatement(statement)
	if err != nil {
		return CreateTransactionResponse{}, err
	}
	return CreateTransactionResponse{
		Limit:   clientLimit,
		Balance: statement.Summary.Total,
	}, nil
}

func futureValueLessThanLimit(value int64, currentBalance int64, clientLimit int64) bool {
	newValue := currentBalance - value
	if newValue < 0 {
		return -newValue > clientLimit
	}
	return false
}
