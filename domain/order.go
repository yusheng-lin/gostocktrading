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
	OrderID  int
	Symbol   string
	Price    float64
	Quantity int
	Type     OrderType
	Time     time.Time
}

// PriorityQueue is a min-heap data structure for orders based on their timestamps.
type PriorityQueue []Order

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Price == pq[j].Price {
		// Orders with the same price are ordered by their timestamps (FIFO)
		return pq[i].Time.Before(pq[j].Time)
	}
	return pq[i].Price > pq[j].Price // For buy orders, higher price is preferred
}

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(Order)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
