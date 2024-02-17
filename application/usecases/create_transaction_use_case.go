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

func NewCreateTransactionUseCase(repository repositories.ClientRepository,
	lockManager lock.LockManager) CreateTransactionUseCase {
	return CreateTransactionUseCase{
		repository:  repository,
		lockManager: lockManager,
	}
}

func (useCase CreateTransactionUseCase) Execute(
	clientId int,
	value int64,
	transactionType string,
	description string,
	clientLimit int64) error {
	lockAcquiringAttempt := 0
	for lockAcquiringAttempt < 3 {
		err := useCase.lockManager.Acquire(strconv.Itoa(clientId))
		if err != nil && err.Error() == "LOCK_ALREADY_ACQUIRED" {
			lockAcquiringAttempt++
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if lockAcquiringAttempt == 3 {
		return fmt.Errorf("LOCK_NOT_ALLOWED")
	}

	currentBalance, err := useCase.repository.GetBalance(clientId)
	if err != nil {
		err = useCase.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			return err
		}
		return err
	}
	if transactionType == "d" && futureValueLessThanLimit(value, currentBalance, clientLimit) {
		err = useCase.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			return err
		}
		return fmt.Errorf("LIMIT_NOT_ALLOWED")
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
	err = useCase.repository.SaveTransaction(transaction)
	if err != nil {
		err = useCase.lockManager.Release(strconv.Itoa(clientId))
		if err != nil {
			return err
		}
		return err
	}
	err = useCase.lockManager.Release(strconv.Itoa(clientId))
	if err != nil {
		return err
	}
	return nil
}

func futureValueLessThanLimit(value int64, currentBalance int64, clientLimit int64) bool {
	newValue := currentBalance - value
	if newValue < 0 {
		return -newValue > clientLimit
	}
	return false
}
