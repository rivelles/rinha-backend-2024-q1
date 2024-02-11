package repositories

import "github.com/rivelles/rinha-backend-2024-q1/model"

type ClientRepository interface {
	SaveTransaction(transaction model.Transaction)
	GetStatement(clientId int) model.Statement
}