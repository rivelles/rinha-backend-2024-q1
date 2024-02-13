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

func NewScyllaRepository() *ScyllaRepository {
	cfg := *gocql.NewCluster("localhost:9042")
	session, _ := gocql.NewSession(cfg)
	session.Query("create keyspace if not exists rinha with replication = {'class':'SimpleStrategy', 'replication_factor':1}").Exec()

	session.Query("create table if not exists rinha.transactions(client_id int, value int, type text, description text, PRIMARY KEY(client_id))").Exec()

	fmt.Println("Created keyspace")

	return &ScyllaRepository{session: session}
}

func (s ScyllaRepository) SaveTransaction(transaction model.Transaction) {
	query := fmt.Sprintf("insert into rinha.transactions (client_id, value, type, description) VALUES "+
		"(%v, %v, '%v', '%v')", transaction.ClientId, transaction.Value, transaction.TransactionType, transaction.Description)
	s.session.Query(query).Exec()
}

func (s ScyllaRepository) GetStatement(clientId int) model.Statement {
	query := fmt.Sprintf("select * from rinha.transactions WHERE client_id = %v", clientId)
	scanner := s.session.Query(query).Iter().Scanner()
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
