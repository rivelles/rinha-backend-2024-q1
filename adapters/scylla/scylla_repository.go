package scylla

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/rivelles/rinha-backend-2024-q1/model"
	"log"
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
	scanner := s.session.Query(`SELECT * FROM transactions WHERE client_id = ?`,
		clientId).Iter().Scanner()
	var transactions []model.Transaction
	for scanner.Next() {
		var (
			value           int
			transactionType string
			description     string
		)
		err := scanner.Scan(&value, &transactionType, &description)
		if err != nil {
			log.Fatal(err)
		}
		i := append(transactions, model.Transaction{
			clientId,
			value,
			transactionType,
			description,
		})
		transactions = i
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return model.Statement{ClientId: clientId, Transactions: transactions}
}
