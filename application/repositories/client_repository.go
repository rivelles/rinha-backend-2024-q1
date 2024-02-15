package repositories

import "github.com/rivelles/rinha-backend-2024-q1/model"

type ClientRepository interface {
	SaveTransaction(transaction model.Transaction) error
	GetStatement(clientId int, clientLimit int) (model.Statement, error)
	GetBalance(clientId int) (int, error)
}
