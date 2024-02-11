package main

import (
	"fmt"
	"github.com/rivelles/rinha-backend-2024-q1/application/usecases"
	"net/http"
)

type App struct {
	createTransactionUseCase usecases.CreateTransactionUseCase,
	getStatementUseCase usecases.GetStatementUseCase,
}

func runServer() {
	http.HandleFunc("/clientes/{id}/transacoes", HandleCreateTransaction)
	http.HandleFunc("/clientes/{id}/extrato", HandleGetStatement)

	http.ListenAndServe(":9091", nil)
}

func HandleCreateTransaction(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Received a request to create transaction")
}

func HandleGetStatement(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Received a request to get statement")
}

func main() {
	runServer()
}
