package model

import (
	"go/types"
)

type Extrato struct {
	idCliente  int
	transacoes types.Slice
}
