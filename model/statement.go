package model

type Statement struct {
	Summary      Summary       `json:"saldo"`
	Transactions []Transaction `json:"ultimas_transacoes"`
}

type Summary struct {
	Total       int    `json:"total"`
	GeneratedAt string `json:"data_extrato"`
	Limit       int    `json:"limite"`
}
