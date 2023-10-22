package service

import (
	"context"
	"sync"
	"time"
	dm "yusheng/stocktrading/domain"
)

type Engine struct {
	Cache map[string]*dm.OrderBook
	lock  sync.RWMutex
}

func NewEngine() dm.IEngine {
	return &Engine{
		Cache: make(map[string]*dm.OrderBook),
	}
}

func (engine *Engine) AddOrder(order *dm.Order) {
	book := engine.GetOrderBook(order.Symbol)
	book.AddOrder(order)
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
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			select {
			case <-ctx.Done():
				return // kill goroutine
			default:
				engine.lock.RLock()
				for _, book := range engine.Cache {
					book.MatchOrders()
				}
				engine.lock.RUnlock()
			}
		}
	}()
	return cancel
}
