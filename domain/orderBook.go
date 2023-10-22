package domain

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type OrderBook struct {
	BuyOrders  PriorityQueue
	SellOrders PriorityQueue
	lock       sync.Mutex
}

func NewOrderBook() *OrderBook {
	book := &OrderBook{
		BuyOrders:  make(PriorityQueue, 0),
		SellOrders: make(PriorityQueue, 0),
	}
	heap.Init(&book.BuyOrders)
	heap.Init(&book.SellOrders)
	return book
}

type Match struct {
	BuyOrder  Order
	SellOrder Order
	Quantity  int
	Price     float64
}

func (orderBook *OrderBook) AddOrder(trader, symbol string, price float64, quantity int, orderType OrderType) {
	orderBook.lock.Lock()
	defer orderBook.lock.Unlock()
	order := Order{
		OrderID:  len(orderBook.BuyOrders) + len(orderBook.SellOrders) + 1,
		Symbol:   symbol,
		Price:    price,
		Quantity: quantity,
		Type:     orderType,
		Time:     time.Now(),
	}
	if orderType == Buy {
		heap.Push(&orderBook.BuyOrders, order)
		fmt.Printf("%s placed a buy order for %s at $%.2f\n", trader, symbol, price)
	} else {
		heap.Push(&orderBook.SellOrders, order)
		fmt.Printf("%s placed a sell order for %s at $%.2f\n", trader, symbol, price)
	}
}

// matchOrders matches buy and sell orders in the order book and returns the list of matched trades.
func (orderBook *OrderBook) MatchOrders() []Match {
	orderBook.lock.Lock()
	defer orderBook.lock.Unlock()

	var matches []Match
	for len(orderBook.BuyOrders) > 0 && len(orderBook.SellOrders) > 0 {
		buyOrder := heap.Pop(&orderBook.BuyOrders).(Order)
		sellOrder := heap.Pop(&orderBook.SellOrders).(Order)

		if buyOrder.Symbol == sellOrder.Symbol && buyOrder.Price >= sellOrder.Price {
			quantity := min(buyOrder.Quantity, sellOrder.Quantity)
			price := sellOrder.Price
			match := Match{
				BuyOrder:  buyOrder,
				SellOrder: sellOrder,
				Quantity:  quantity,
				Price:     price,
			}
			matches = append(matches, match)

			buyOrder.Quantity -= quantity
			sellOrder.Quantity -= quantity

			if buyOrder.Quantity > 0 {
				heap.Push(&orderBook.BuyOrders, buyOrder)
			}
			if sellOrder.Quantity > 0 {
				heap.Push(&orderBook.SellOrders, sellOrder)
			}
		}
	}
	return matches
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
