package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rinha-backend-2024-q1/application/repositories"
	"rinha-backend-2024-q1/application/usecases"
	"strconv"
)

type App struct {
	createTransactionUseCase usecases.CreateTransactionUseCase
	getStatementUseCase      usecases.GetStatementUseCase
}

type TransactionRequest struct {
	Value int64  `json:"valor"`
	Type  string `json:"tipo"`
	Desc  string `json:"descricao"`
}

// Roubando :)
var limitByClientId = map[int]int64{
	1: 100000,
	2: 80000,
	3: 1000000,
	4: 10000000,
	5: 500000,
}

func NewApp(repository repositories.ClientRepository) *App {
	return &App{
		usecases.NewCreateTransactionUseCase(repository),
		usecases.NewGetStatementUseCase(repository),
	}
}

func (a App) Run(port int) {
	http.HandleFunc("/clientes/{id}/transacoes", a.HandleCreateTransaction)
	http.HandleFunc("/clientes/{id}/extrato", a.HandleGetStatement)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (a App) HandleCreateTransaction(writer http.ResponseWriter, req *http.Request) {
	clientId, _ := strconv.Atoi(req.PathValue("id"))
	var transactionRequest TransactionRequest
	err := json.NewDecoder(req.Body).Decode(&transactionRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	validationErr := validateTransactionRequest(transactionRequest)
	if validationErr != nil {
		http.Error(writer, validationErr.Error(), http.StatusUnprocessableEntity)
		return
	}
	response, err := a.createTransactionUseCase.Execute(
		clientId,
		transactionRequest.Value,
		transactionRequest.Type,
		transactionRequest.Desc,
		limitByClientId[clientId],
	)
	if err != nil {
		if err.Error() == "LIMIT_NOT_ALLOWED" {
			http.Error(writer, "", http.StatusUnprocessableEntity)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	resultJson, _ := json.Marshal(response)
	writer.Header().Set("Content-Type", "text/json; charset=utf-8")
	writer.Write(resultJson)
}

func validateTransactionRequest(request TransactionRequest) error {
	if request.Type != "d" && request.Type != "c" {
		return fmt.Errorf("TYPE_NOT_ALLOWED")
	}
	if len(request.Desc) == 0 || len(request.Desc) > 10 {
		return fmt.Errorf("DESC_NOT_ALLOWED")
	}
	if request.Value <= 0 {
		return fmt.Errorf("VALUE_NOT_ALLOWED")
	}
	return nil
}

func (a App) HandleGetStatement(writer http.ResponseWriter, req *http.Request) {
	clientId, _ := strconv.Atoi(req.PathValue("id"))
	if clientId > 5 || clientId < 0 {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		statement, err := a.getStatementUseCase.Execute(clientId)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "text/json; charset=utf-8")
		writer.WriteHeader(200)
		statementJson, _ := json.Marshal(statement)
		writer.Write(statementJson)
	}
}
