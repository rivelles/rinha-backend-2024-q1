package main

import (
	"rinha-backend-2024-q1/adapters/scylla"
)

func main() {
	repository := scylla.NewScyllaRepository()
	app := NewApp(repository)

	app.Run(9091)
}
