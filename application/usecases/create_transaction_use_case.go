package usecases

import (
	"fmt"
	"rinha-backend-2024-q1/application/lock"
	"rinha-backend-2024-q1/application/repositories"
	"rinha-backend-2024-q1/model"
	"strconv"
	"time"
)

type CreateTransactionUseCase struct {
	repository  repositories.ClientRepository
	lockManager lock.LockManager
}

type CreateTransactionResponse struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

func NewCreateTransactionUseCase(repository repositories.ClientRepository,
	lockManager lock.LockManager) CreateTransactionUseCase {
	return CreateTransactionUseCase{
		repository:  repository,
		lockManager: lockManager,
	}
}

func (u CreateTransactionUseCase) Execute(
	clientId int,
	value int64,
	transactionType string,
	description string,
	clientLimit int64) (CreateTransactionResponse, error) {
	lockAcquired := false
	for lockAcquired == false {
		err := u.lockManager.Acquire(strconv.Itoa(clientId))
		if err != nil && err.Error() == "LOCK_ALREADY_ACQUIRED" {
			println("Lock already acquired, will try again...")
			millisToWait := time.Duration(10)
			time.Sleep(millisToWait * time.Millisecond)
			continue
		}
		lockAcquired = true
	}

	currentBalance, err := u.repository.GetBalance(clientId)
	if err != nil {
		err = u.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			println("Error: Failed to release lock after get balance")
			return CreateTransactionResponse{}, err
		}
		println("Error: Failed to get balance")
		return CreateTransactionResponse{}, err
	}
	if transactionType == "d" && futureValueLessThanLimit(value, currentBalance, clientLimit) {
		err = u.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			println("Error: Failed to release lock before returning limit not allowed error")
			return CreateTransactionResponse{}, err
		}
		return CreateTransactionResponse{}, fmt.Errorf("LIMIT_NOT_ALLOWED")
	}
	newBalance := currentBalance
	if transactionType == "d" {
		newBalance -= value
	} else {
		newBalance += value
	}
	transaction := model.Transaction{
		ClientId:        clientId,
		Timestamp:       time.Now().UnixMilli(),
		Value:           value,
		CurrentBalance:  newBalance,
		TransactionType: transactionType,
		Description:     description,
	}
	err = u.repository.SaveTransaction(transaction)
	if err != nil {
		err = u.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			println("Error: Failed to release lock after saving transaction")
			return CreateTransactionResponse{}, err
		}
		println("Error: Failed to save transaction")
		return CreateTransactionResponse{}, err
	}
	response := CreateTransactionResponse{
		Limit:   clientLimit,
		Balance: newBalance,
	}
	u.lockManager.Release(strconv.Itoa(clientId))
	return response, nil
}

func futureValueLessThanLimit(value int64, currentBalance int64, clientLimit int64) bool {
	newValue := currentBalance - value
	if newValue < 0 {
		return -newValue > clientLimit
	}
	return false
}
