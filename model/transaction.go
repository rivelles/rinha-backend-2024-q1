package model

type Transaction struct {
	Value           int64  `json:"valor"`
	TransactionType string `json:"tipo"`
	Description     string `json:"descricao"`
	Timestamp       string `json:"realizada_em"`
}
