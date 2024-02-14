package scylla

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/rivelles/rinha-backend-2024-q1/model"
	"log"
	"strconv"
	"time"
)

type ScyllaRepository struct {
	session *gocql.Session
}

func NewScyllaRepository() *ScyllaRepository {
	cfg := *gocql.NewCluster("localhost:9042")
	session, _ := gocql.NewSession(cfg)
	session.Query("create keyspace if not exists rinha with replication = {'class':'SimpleStrategy', 'replication_factor':1}").Exec()

	session.Query("create table if not exists rinha.transactions(client_id int, value int, type text, description text, timestamp bigint, PRIMARY KEY(client_id, timestamp))").Exec()

	fmt.Println("Created keyspace")

	return &ScyllaRepository{session: session}
}

func (s ScyllaRepository) SaveTransaction(transaction model.Transaction) {
	query := fmt.Sprintf("insert into rinha.transactions (client_id, value, type, description, timestamp) VALUES "+
		"(%v, %v, '%v', '%v', %v)", transaction.ClientId, transaction.Value, transaction.TransactionType, transaction.Description, transaction.Timestamp)
	s.session.Query(query).Exec()
}

func (s ScyllaRepository) GetStatement(clientId int, clientLimit int) model.Statement {
	query := fmt.Sprintf("select value, type, description, timestamp "+
		"from rinha.transactions "+
		"WHERE client_id = %v "+
		"ORDER BY timestamp DESC "+
		"PER PARTITION LIMIT 10", clientId)

	scanner := s.session.Query(query).Iter().Scanner()
	var transactions []model.Transaction
	balance := 0
	for scanner.Next() {
		var (
			value           string
			transactionType string
			description     string
			timestamp       string
		)
		err := scanner.Scan(&value, &transactionType, &description, &timestamp)
		if err != nil {
			log.Fatal(err)
		}
		intValue, _ := strconv.Atoi(value)
		timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
		i := append(transactions, model.Transaction{
			ClientId:        clientId,
			Value:           intValue,
			TransactionType: transactionType,
			Description:     description,
			Timestamp:       timestampInt,
		})
		transactions = i
		if transactionType == "c" {
			balance += intValue
		} else {
			balance -= intValue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	summary := model.Summary{
		Total:       balance,
		GeneratedAt: time.Now().String(),
		Limit:       clientLimit,
	}
	return model.Statement{
		Summary:      summary,
		Transactions: transactions,
	}
}

func (s ScyllaRepository) GetBalance(clientId int) int {
	query := fmt.Sprintf("select value, type "+
		"from rinha.transactions "+
		"where client_id = %v", clientId)
	scanner := s.session.Query(query).Iter().Scanner()
	balance := 0
	for scanner.Next() {
		var (
			value           string
			transactionType string
		)
		err := scanner.Scan(&value, &transactionType)
		if err != nil {
			log.Fatal(err)
		}
		intValue, _ := strconv.Atoi(value)
		if transactionType == "c" {
			balance += intValue
		} else {
			balance -= intValue
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return balance
}
