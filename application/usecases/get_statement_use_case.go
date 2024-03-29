package usecases

import (
	"rinha-backend-2024-q1/application/repositories"
	"rinha-backend-2024-q1/model"
)

type GetStatementUseCase struct {
	repository repositories.ClientRepository
}

func NewGetStatementUseCase(repository repositories.ClientRepository) GetStatementUseCase {
	return GetStatementUseCase{repository}
}

func (g GetStatementUseCase) Execute(clientId int) (model.Statement, error) {
	statement, err := g.repository.GetStatement(clientId)
	if err != nil {
		return model.Statement{}, err
	}
	return statement, nil
}
