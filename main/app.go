package main

import (
	"github.com/rivelles/rinha-backend-2024-q1/application/usecases"
	"net/http"
	"strconv"
)

type App struct {
	createTransactionUseCase usecases.CreateTransactionUseCase
	getStatementUseCase      usecases.GetStatementUseCase
}

func NewApp() *App {
	return &App{
		usecases.NewCreateTransactionUseCase(),
		usecases.NewGetStatementUseCase(),
	}
}

func (a App) Run(port int) {
	http.HandleFunc("/clientes/{id}/transacoes", HandleCreateTransaction)
	http.HandleFunc("/clientes/{id}/extrato", HandleGetStatement)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func HandleCreateTransaction(writer http.ResponseWriter, req *http.Request) {

}

func HandleGetStatement(writer http.ResponseWriter, req *http.Request) {

}
