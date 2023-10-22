package domain

import "context"

type IEngine interface {
	AddOrder(order *Order)
	MatchOrders() context.CancelFunc
}
