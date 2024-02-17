package model

type Transaction struct {
	ClientId        int    `json:"-"`
	Value           int64  `json:"valor"`
	CurrentBalance  int64  `json:"-"`
	TransactionType string `json:"tipo"`
	Description     string `json:"descricao"`
	Timestamp       int64  `json:"-"`
	TimestampStr    string `json:"realizada_em"`
}
