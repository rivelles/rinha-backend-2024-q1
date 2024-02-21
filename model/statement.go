package model

type Statement struct {
	ClientId          int           `json:"-"`
	TotalTransactions int64         `json:"-"`
	Summary           Summary       `json:"saldo"`
	Transactions      []Transaction `json:"ultimas_transacoes"`
}

type Summary struct {
	Total       int64  `json:"total"`
	GeneratedAt string `json:"data_extrato"`
	Limit       int64  `json:"limite"`
}
