package main

import (
	"fmt"
	"sync"
	dm "yusheng/stocktrading/domain"
)

func main() {
	orderBook := dm.NewOrderBook()
	var wg sync.WaitGroup
	traderCount := 6
	traderSymbols := []string{"AAPL", "GOOG", "TSLA"} // Different symbols for traders

	for i := 0; i < traderCount; i++ {
		trader := fmt.Sprintf("Trader%d", i+1)
		symbol := traderSymbols[i%len(traderSymbols)]
		wg.Add(1)
		go func(trader, symbol string) {
			defer wg.Done()
			orderBook.AddOrder(trader, symbol, 150.0, 50, dm.Buy)
			orderBook.AddOrder(trader, symbol, 150.0, 30, dm.Sell)
		}(trader, symbol)
	}
	wg.Wait()
	matches := orderBook.MatchOrders()

	for _, match := range matches {
		fmt.Printf("Matched - Buy Order %d with Sell Order %d:\n", match.BuyOrder.OrderID, match.SellOrder.OrderID)
		fmt.Printf("Symbol: %s, Price: %.2f, Quantity: %d\n", match.BuyOrder.Symbol, match.Price, match.Quantity)
		fmt.Println()
	}
}
