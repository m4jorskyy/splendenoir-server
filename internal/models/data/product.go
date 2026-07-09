package data

import (
	"database/sql"
)

type Product struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Material  string       `json:"material"`
	Fineness  string       `json:"fineness"`
	Length    float64      `json:"length"`
	Width     float64      `json:"width"`
	Price     float64      `json:"price"`
	Quantity  int64        `json:"quantity"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
