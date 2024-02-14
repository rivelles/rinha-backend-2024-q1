package usecases

import (
	"fmt"
	"github.com/rivelles/rinha-backend-2024-q1/application/lock"
	"github.com/rivelles/rinha-backend-2024-q1/application/repositories"
	"github.com/rivelles/rinha-backend-2024-q1/model"
	"time"
)

type CreateTransactionUseCase struct {
	repository  repositories.ClientRepository
	lockManager lock.LockManager
}

func NewCreateTransactionUseCase(repository repositories.ClientRepository, lockManager lock.LockManager) CreateTransactionUseCase {
	return CreateTransactionUseCase{
		repository:  repository,
		lockManager: lockManager,
	}
}

func (useCase CreateTransactionUseCase) Execute(clientId int, value int, transactionType string, description string, clientLimit int) error {
	currentBalance := useCase.repository.GetBalance(clientId)
	if transactionType == "d" && futureValueLessThanLimit(clientId, value, currentBalance, clientLimit) {
		return fmt.Errorf("transaction not allowed: future balance would be less than limit")
	}
	transaction := model.Transaction{
		ClientId:        clientId,
		Value:           value,
		TransactionType: transactionType,
		Description:     description,
		Timestamp:       time.Now().UnixMilli(),
	}
	useCase.repository.SaveTransaction(transaction)
	return nil
}

func futureValueLessThanLimit(clientId int, value int, currentBalance int, clientLimit int) bool {
	return currentBalance-value < clientLimit
}
