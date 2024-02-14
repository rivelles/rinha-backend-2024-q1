package model

type Transaction struct {
	ClientId        int    `json:"-"`
	Value           int    `json:"valor"`
	CurrentBalance  int    `json:"-"`
	TransactionType string `json:"tipo"`
	Description     string `json:"descricao"`
	Timestamp       int64  `json:"realizada_em"`
}
