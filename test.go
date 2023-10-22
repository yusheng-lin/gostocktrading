package main

import (
	"sync"
	"time"
	dm "yusheng/stocktrading/domain"
)

var orderId int = 0

func placeOrders_test(engine dm.IEngine, wg *sync.WaitGroup) {
	for _, custome := range dummy1() {
		addOrders(engine, custome, wg)
	}
}

func addOrders(engine dm.IEngine, orders []*dm.Order, wg *sync.WaitGroup) {
	wg.Add(1)
	go func(engine dm.IEngine, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, order := range orders {
			engine.AddOrder(order)
		}
	}(engine, wg)
}

func getOrderId() int {
	orderId++
	return orderId
}

// customer1(c1)
// 1. buy AAPL $50 Qty:30
// 2. buy GOOG $35 Qty:20
// 3. sell TSLA $40 Qty:20
// customer2(c2)
// 1. sell AAPL $50 Qty:40
// 2. sell GOOG $30 Qty:10
// 3. buy  TSLA $45 Qty:30
// customer3(c3)
// 1. sell TSLA $45 Qty:20
// 2. buy  AAPL $50 Qty:10
// 3. sell GOOG $35 Qty:5
// expect
// c1 buy 30 $50 AAPL from c2, 10 $30 GOOG from c2, 5 $35 GOOG from c3
// c2 buy 20 $40 TSLA from c1, 10 $45 TSLA from c3
// c3 buy 10 $50 AAPL from c2
func dummy1() [][]*dm.Order {
	customer1 := []*dm.Order{
		{
			OrderID:    getOrderId(),
			Symbol:     "AAPL",
			Price:      50,
			Quantity:   30,
			Type:       dm.Buy,
			Time:       time.Now(),
			CustomerID: 1,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "GOOG",
			Price:      35,
			Quantity:   20,
			Type:       dm.Buy,
			Time:       time.Now(),
			CustomerID: 1,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "TSLA",
			Price:      40,
			Quantity:   20,
			Type:       dm.Sell,
			Time:       time.Now(),
			CustomerID: 1,
		},
	}
	customer2 := []*dm.Order{
		{
			OrderID:    getOrderId(),
			Symbol:     "AAPL",
			Price:      50,
			Quantity:   40,
			Type:       dm.Sell,
			Time:       time.Now(),
			CustomerID: 2,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "GOOG",
			Price:      30,
			Quantity:   10,
			Type:       dm.Sell,
			Time:       time.Now(),
			CustomerID: 2,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "TSLA",
			Price:      45,
			Quantity:   30,
			Type:       dm.Buy,
			Time:       time.Now(),
			CustomerID: 2,
		},
	}
	customer3 := []*dm.Order{
		{
			OrderID:    getOrderId(),
			Symbol:     "TSLA",
			Price:      45,
			Quantity:   30,
			Type:       dm.Sell,
			Time:       time.Now(),
			CustomerID: 3,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "AAPL",
			Price:      50,
			Quantity:   10,
			Type:       dm.Buy,
			Time:       time.Now(),
			CustomerID: 3,
		},
		{
			OrderID:    getOrderId(),
			Symbol:     "GOOG",
			Price:      35,
			Quantity:   5,
			Type:       dm.Sell,
			Time:       time.Now(),
			CustomerID: 3,
		},
	}
	return [][]*dm.Order{customer1, customer2, customer3}
}
