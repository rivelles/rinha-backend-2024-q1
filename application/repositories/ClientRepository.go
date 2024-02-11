package repositories

import "github.com/rivelles/rinha-backend-2024-q1/model"

type StatementRepository interface {
	Save(transaction model.Transaction)
	GetStatement(clientId int)
}
