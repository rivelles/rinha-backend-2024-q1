package usecases

import (
	"github.com/rivelles/rinha-backend-2024-q1/application/repositories"
	"github.com/rivelles/rinha-backend-2024-q1/model"
)

type CreateTransactionUseCase struct {
	repository repositories.ClientRepository
}

func (useCase CreateTransactionUseCase) execute(
	clientId int,
	value int,
	transactionType string,
	description string) {
	transaction := model.Transaction{
		ClientId:        clientId,
		Value:           value,
		TransactionType: transactionType,
		Description:     description,
	}
	useCase.repository.SaveTransaction(transaction)
}
