package main

import (
	"encoding/json"
	"github.com/rivelles/rinha-backend-2024-q1/application/lock"
	"github.com/rivelles/rinha-backend-2024-q1/application/repositories"
	"github.com/rivelles/rinha-backend-2024-q1/application/usecases"
	"net/http"
	"strconv"
)

type App struct {
	createTransactionUseCase usecases.CreateTransactionUseCase
	getStatementUseCase      usecases.GetStatementUseCase
}

// Roubando :)
var limitByClientId = map[int]int{
	1: 100000,
	2: 80000,
	3: 1000000,
	4: 10000000,
	5: 500000,
}

func NewApp(repository repositories.ClientRepository, lockManager lock.LockManager) *App {
	return &App{
		usecases.NewCreateTransactionUseCase(repository, lockManager),
		usecases.NewGetStatementUseCase(repository),
	}
}

func (a App) Run(port int) {
	http.HandleFunc("/clientes/{id}/transacoes", a.HandleCreateTransaction)
	http.HandleFunc("/clientes/{id}/extrato", a.HandleGetStatement)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (a App) HandleCreateTransaction(writer http.ResponseWriter, req *http.Request) {
	a.createTransactionUseCase.Execute(1, 1, "c", "aaa", limitByClientId[1])
}

func (a App) HandleGetStatement(writer http.ResponseWriter, req *http.Request) {
	clientId, _ := strconv.Atoi(req.PathValue("id"))
	if clientId > 5 || clientId < 0 {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		statement := a.getStatementUseCase.Execute(clientId, limitByClientId[clientId])
		writer.Header().Set("Content-Type", "text/json; charset=utf-8")
		writer.WriteHeader(200)
		statementJson, _ := json.Marshal(statement)
		writer.Write(statementJson)
	}
}
