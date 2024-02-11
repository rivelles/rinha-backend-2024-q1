package model

import (
	"go/types"
)

type Statement struct {
	clientId     int
	transactions types.Slice
}
