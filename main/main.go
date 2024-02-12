package main

import (
	"github.com/rivelles/rinha-backend-2024-q1/adapters/scylla"
)

func main() {
	repository := scylla.NewScyllaRepository()
	app := NewApp(repository, nil)

	app.Run(9091)
}
