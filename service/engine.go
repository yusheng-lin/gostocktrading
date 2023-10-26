package service

import (
	"context"
	"fmt"
	"sync"
	dm "yusheng/stocktrading/domain"
)

type Engine struct {
	Cache map[string]*dm.OrderBook
	lock  sync.RWMutex
	ch    chan string
}

func NewEngine() dm.IEngine {
	return &Engine{
		Cache: make(map[string]*dm.OrderBook),
		ch:    make(chan string),
	}
}

func (engine *Engine) AddOrder(order *dm.Order) {
	book := engine.GetOrderBook(order.Symbol)
	book.AddOrder(order)
	engine.ch <- order.Symbol
}

func (engine *Engine) GetOrderBook(symbol string) *dm.OrderBook {
	engine.lock.RLock()
	book, ok := engine.Cache[symbol]
	engine.lock.RUnlock()
	if !ok {
		engine.lock.Lock()
		book, ok = engine.Cache[symbol]
		if !ok {
			book = dm.NewOrderBook(symbol)
			engine.Cache[symbol] = book
		}
		engine.lock.Unlock()
	}
	return book
}

func (engine *Engine) MatchOrders() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("init 4 workers")
	for i := 0; i < 4; i++ {
		go func(engine *Engine, ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("recieved cancel signal")
					return // kill goroutine
				case symbol := <-engine.ch:
					engine.lock.RLock()
					engine.GetOrderBook(symbol).MatchOrders()
					engine.lock.RUnlock()
				default:
				}
			}
		}(engine, ctx)
	}
	return cancel
}
