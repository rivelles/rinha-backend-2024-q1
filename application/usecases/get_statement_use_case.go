package usecases

import (
	"github.com/rivelles/rinha-backend-2024-q1/application/repositories"
	"github.com/rivelles/rinha-backend-2024-q1/model"
)

type GetStatementUseCase struct {
	repository repositories.ClientRepository
}

func NewGetStatementUseCase(repository repositories.ClientRepository) GetStatementUseCase {
	return GetStatementUseCase{repository}
}

func (g GetStatementUseCase) Execute(clientId int, clientLimit int) model.Statement {
	return g.repository.GetStatement(clientId, clientLimit)
}
