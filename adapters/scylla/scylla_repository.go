package scylla

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"rinha-backend-2024-q1/model"
	"strconv"
	"time"
)

type ScyllaRepository struct {
	session *gocql.Session
}

func NewScyllaRepository() *ScyllaRepository {
	host, exists := os.LookupEnv("SCYLLA_HOST")
	if !exists {
		host = "localhost"
	}
	port, exists := os.LookupEnv("SCYLLA_PORT")
	if !exists {
		port = "9042"
	}
	cfg := *gocql.NewCluster(fmt.Sprintf("%v:%v", host, port))
	session, _ := gocql.NewSession(cfg)
	session.Query("create keyspace if not exists rinha with replication = {'class':'SimpleStrategy', 'replication_factor':1}").Exec()

	session.Query("create table if not exists rinha.transactions" +
		"(client_id int, " +
		"timestamp bigint, " +
		"value bigint, " +
		"current_balance bigint, " +
		"type text, " +
		"description text, " +
		"PRIMARY KEY(client_id, timestamp)" +
		")").Exec()

	fmt.Println("Created keyspace")

	return &ScyllaRepository{session: session}
}

func (s ScyllaRepository) SaveTransaction(transaction model.Transaction) error {
	query := fmt.Sprintf("insert into rinha.transactions (client_id, value, current_balance, type, description, timestamp) VALUES "+
		"(%v, %v, %v, '%v', '%v', %v)", transaction.ClientId, transaction.Value, transaction.CurrentBalance, transaction.TransactionType, transaction.Description, transaction.Timestamp)
	err := s.session.Query(query).Exec()
	return err
}

func (s ScyllaRepository) GetStatement(clientId int, clientLimit int64) (model.Statement, error) {
	query := fmt.Sprintf("select value, current_balance, type, description, timestamp "+
		"from rinha.transactions "+
		"WHERE client_id = %v "+
		"ORDER BY timestamp DESC "+
		"PER PARTITION LIMIT 10", clientId)

	scanner := s.session.Query(query).Iter().Scanner()
	var balance int64 = 0
	firstRow := true
	var transactions []model.Transaction
	for scanner.Next() {
		var (
			value           int64
			currentBalance  int64
			transactionType string
			description     string
			timestamp       string
		)
		err := scanner.Scan(&value, &currentBalance, &transactionType, &description, &timestamp)
		if err != nil {
			log.Fatal(err)
		}
		if firstRow {
			balance = currentBalance
			firstRow = false
		}
		timestampInt, _ := strconv.ParseInt(timestamp, 10, 64)
		i := append(transactions, model.Transaction{
			ClientId:        clientId,
			Value:           value,
			TransactionType: transactionType,
			Description:     description,
			Timestamp:       timestampInt,
			TimestampStr:    time.UnixMilli(timestampInt).Format(time.RFC3339),
		})
		transactions = i
	}
	if err := scanner.Err(); err != nil {
		return model.Statement{}, err
	}

	summary := model.Summary{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Total:       balance,
		Limit:       clientLimit,
	}
	return model.Statement{
		Summary:      summary,
		Transactions: transactions,
	}, nil
}

func (s ScyllaRepository) GetBalance(clientId int) (int64, error) {
	query := fmt.Sprintf("SELECT current_balance "+
		"FROM rinha.transactions "+
		"WHERE client_id = %v "+
		"ORDER BY timestamp DESC "+
		"PER PARTITION LIMIT 1", clientId)
	scanner := s.session.Query(query).Iter().Scanner()
	var balance int64 = 0
	for scanner.Next() {
		var (
			value string
		)
		err := scanner.Scan(&value)
		if err != nil {
			log.Fatal(err)
		}
		balance, _ = strconv.ParseInt(value, 10, 64)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return balance, nil
}
