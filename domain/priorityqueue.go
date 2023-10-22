package domain

// PriorityQueue is a min-heap data structure for orders based on their timestamps.
type PriorityQueue struct {
	orders    []Order
	orderType OrderType
}

func NewPriorityQueue(orderType OrderType) *PriorityQueue {
	return &PriorityQueue{
		orders:    make([]Order, 0),
		orderType: orderType,
	}
}

func (pq PriorityQueue) Len() int { return len(pq.orders) }

func (pq PriorityQueue) Less(i, j int) bool {
	q := pq.orders
	if q[i].Price == q[j].Price {
		// Orders with the same price are ordered by their timestamps (FIFO)
		return q[i].Time.Before(q[j].Time)
	}
	if pq.orderType == Buy {
		return q[i].Price > q[j].Price
	}
	return q[i].Price < q[j].Price
}

func (pq PriorityQueue) Swap(i, j int) {
	q := pq.orders
	q[i], q[j] = q[j], q[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(Order)
	pq.orders = append(pq.orders, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.orders
	n := len(old)
	item := old[n-1]
	pq.orders = old[0 : n-1]
	return item
}
