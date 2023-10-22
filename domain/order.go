package domain

import (
	"time"
)

type OrderType int

const (
	Buy OrderType = iota
	Sell
)

type Order struct {
	OrderID    int
	Symbol     string
	Price      float64
	Quantity   int
	Type       OrderType
	Time       time.Time
	CustomerID int
}
