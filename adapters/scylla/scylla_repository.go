package scylla

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"os"
	"rinha-backend-2024-q1/model"
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
		"(client_id tinyint, " +
		"current_balance bigint, " +
		"last_transactions text, " +
		"number_transactions bigint, " +
		"PRIMARY KEY(client_id)" +
		")").Exec()

	session.Query("insert into rinha.transactions " +
		"(client_id, current_balance, last_transactions, number_transactions)" +
		"VALUES (1, 0, '', 0)").Exec()

	session.Query("insert into rinha.transactions " +
		"(client_id, current_balance, last_transactions, number_transactions)" +
		"VALUES (2, 0, '', 0)").Exec()

	session.Query("insert into rinha.transactions " +
		"(client_id, current_balance, last_transactions, number_transactions)" +
		"VALUES (3, 0, '', 0)").Exec()

	session.Query("insert into rinha.transactions " +
		"(client_id, current_balance, last_transactions, number_transactions)" +
		"VALUES (4, 0, '', 0)").Exec()

	session.Query("insert into rinha.transactions " +
		"(client_id, current_balance, last_transactions, number_transactions)" +
		"VALUES (5, 0, '', 0)").Exec()

	fmt.Println("Created keyspace")

	return &ScyllaRepository{session: session}
}

func (s ScyllaRepository) SaveStatement(statement model.Statement) error {
	marshalTransactions, _ := json.Marshal(statement.Transactions)
	incTransactionNumber := statement.TotalTransactions + 1
	query := fmt.Sprintf("update rinha.transactions "+
		"set current_balance = %v, "+
		"number_transactions = %v, "+
		"last_transactions = '%v' "+
		"where client_id = %v "+
		"if number_transactions = %v", statement.Summary.Total, incTransactionNumber, string(marshalTransactions), statement.ClientId, statement.TotalTransactions)
	m := make(map[string]interface{})
	applied, err := s.session.Query(query).MapScanCAS(m)
	if err != nil {
		return err
	}
	if !applied {
		return fmt.Errorf("CONFLICT")
	}
	return nil
}

func (s ScyllaRepository) GetStatement(clientId int) (model.Statement, error) {
	query := fmt.Sprintf("select current_balance, last_transactions, number_transactions "+
		"from rinha.transactions "+
		"WHERE client_id = %v", clientId)

	scanner := s.session.Query(query).Iter().Scanner()
	var transactions []model.Transaction
	for scanner.Next() {
		var (
			currentBalance     int64
			lastTransactions   string
			numberTransactions int64
		)
		err := scanner.Scan(&currentBalance, &lastTransactions, &numberTransactions)
		if err != nil {
			log.Fatal(err)
		}
		summary := model.Summary{
			Total:       currentBalance,
			GeneratedAt: time.Now().String(),
		}
		if lastTransactions != "" {
			err = json.Unmarshal([]byte(lastTransactions), &transactions)
			if err != nil {
				return model.Statement{}, err
			}
		} else {
			transactions = []model.Transaction{}
		}
		statement := model.Statement{
			ClientId:          clientId,
			Summary:           summary,
			Transactions:      transactions,
			TotalTransactions: numberTransactions,
		}
		return statement, nil
	}
	return model.Statement{}, fmt.Errorf("NOT_FOUND")
}
