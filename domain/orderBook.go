package domain

import (
	"container/heap"
	"fmt"
	"sync"
)

type OrderBook struct {
	BuyOrders  *PriorityQueue
	SellOrders *PriorityQueue
	Symbol     string
	lock       sync.Mutex
}

func NewOrderBook(symbol string) *OrderBook {
	book := &OrderBook{
		BuyOrders:  NewPriorityQueue(Buy),
		SellOrders: NewPriorityQueue(Sell),
		Symbol:     symbol,
	}
	heap.Init(book.BuyOrders)
	heap.Init(book.SellOrders)
	return book
}

func (orderBook *OrderBook) AddOrder(order *Order) {
	orderBook.lock.Lock()
	defer orderBook.lock.Unlock()
	if order.Type == Buy {
		heap.Push(orderBook.BuyOrders, *order)
		fmt.Printf("c%d place a buy order%d for %s at $%.2f qty:%d\n", order.CustomerID, order.OrderID, orderBook.Symbol, order.Price, order.Quantity)
	} else {
		heap.Push(orderBook.SellOrders, *order)
		fmt.Printf("c%d place a sell order%d for %s at $%.2f qty:%d\n", order.CustomerID, order.OrderID, orderBook.Symbol, order.Price, order.Quantity)
	}
}

// matchOrders matches buy and sell orders in the order book and returns the list of matched trades.
func (orderBook *OrderBook) MatchOrders() {
	orderBook.lock.Lock()
	defer orderBook.lock.Unlock()
	for orderBook.BuyOrders.Len() > 0 && orderBook.SellOrders.Len() > 0 {
		buyOrder := heap.Pop(orderBook.BuyOrders).(Order)
		sellOrder := heap.Pop(orderBook.SellOrders).(Order)
		if buyOrder.Price >= sellOrder.Price {
			quantity := min(buyOrder.Quantity, sellOrder.Quantity)
			price := sellOrder.Price
			fmt.Printf("c%d buy %d $%.2f %s from c%d \n", buyOrder.CustomerID, quantity, price, orderBook.Symbol, sellOrder.CustomerID)
			buyOrder.Quantity -= quantity
			sellOrder.Quantity -= quantity

			if buyOrder.Quantity > 0 {
				heap.Push(orderBook.BuyOrders, buyOrder)
			}
			if sellOrder.Quantity > 0 {
				heap.Push(orderBook.SellOrders, sellOrder)
			}
		}
	}
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
