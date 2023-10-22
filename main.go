package main

import (
	"sync"
	"time"
	srv "yusheng/stocktrading/service"
)

func main() {
	var wg sync.WaitGroup
	engine := srv.NewEngine()
	wg.Add(1)
	go func() {
		defer wg.Done()
		cancel := engine.MatchOrders()
		defer cancel()
		time.Sleep(5 * time.Second)
	}()
	placeOrders_test(engine, &wg)
	wg.Wait()
}
