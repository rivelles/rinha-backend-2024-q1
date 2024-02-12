package scylla

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/rivelles/rinha-backend-2024-q1/model"
)

type ScyllaRepository struct {
	session *gocql.Session
}

func NewScyllaRepository() ScyllaRepository {
	cfg := *gocql.NewCluster("http://localhost:7000")
	session, _ := gocql.NewSession(cfg)

	return ScyllaRepository{session: session}
}

func (s ScyllaRepository) SaveTransaction(transaction model.Transaction) {
	fmt.Println("Saving.")
}

func (s ScyllaRepository) GetStatement(clientId int) model.Statement {
	return model.Statement{
		ClientId:     clientId,
		Transactions: make([]model.Transaction, 0),
	}
}
