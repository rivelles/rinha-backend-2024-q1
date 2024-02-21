package repositories

import "rinha-backend-2024-q1/model"

type ClientRepository interface {
	GetStatement(clientId int) (model.Statement, error)
	SaveStatement(statement model.Statement) error
}
