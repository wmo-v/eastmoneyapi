package model

import (
	"log"

	"github.com/shopspring/decimal"
)

type QuoteQueue struct {
	arr    []*Stockquote
	front  int
	rear   int
	maxLen int

	totalPrice   decimal.Decimal
	highestPrice decimal.Decimal
	lowestPrice  decimal.Decimal
}

func NewQueue(len int) *QuoteQueue {
	return &QuoteQueue{
		arr:    make([]*Stockquote, len+1),
		maxLen: len + 1,
	}
}

// Enqueue 添加新元素
func (q *QuoteQueue) Enqueue(data *Stockquote) {
	if q.isFull() {
		q.dequeue()
	}
	curPrice := decimal.NewFromFloat(data.NewestPrice)
	if q.isEmpty() || q.highestPrice.LessThan(curPrice) {
		q.highestPrice = curPrice
	}
	if q.isEmpty() || curPrice.LessThan(q.lowestPrice) {
		q.lowestPrice = curPrice
	}
	q.totalPrice = q.totalPrice.Add(curPrice)
	q.arr[q.rear] = data
	q.rear = (q.rear + 1) % q.maxLen
}

// GetHighestPrice 获取当前队列的最大值
func (q *QuoteQueue) GetHighestPrice() decimal.Decimal {
	return q.highestPrice
}

// GetLowestPrice 获取当前队列的最小值
func (q *QuoteQueue) GetLowestPrice() decimal.Decimal {
	return q.lowestPrice
}
func (q *QuoteQueue) GrtAvgPrice() decimal.Decimal {
	count := (q.rear - q.front + q.maxLen) % q.maxLen
	return q.totalPrice.Div(decimal.NewFromInt(int64(count)))
}
func (q *QuoteQueue) isFull() bool {
	return (q.rear+1)%q.maxLen == q.front
}

func (q *QuoteQueue) isEmpty() bool {
	return q.front == q.rear
}

func (q *QuoteQueue) dequeue() {
	value := decimal.NewFromFloat(q.arr[q.front].NewestPrice)
	q.front = (q.front + 1) % q.maxLen
	q.totalPrice = q.totalPrice.Sub(value)
	if q.isEmpty() {
		return
	}
	if value.Equal(q.highestPrice) || value.Equal(q.lowestPrice) {
		q.reSelectMinAndMax()
	}
}
func (q *QuoteQueue) reSelectMinAndMax() {
	var max = decimal.NewFromFloat(q.arr[q.front].NewestPrice)
	var min = decimal.NewFromFloat(q.arr[q.front].NewestPrice)
	var idx = q.front
	for idx != q.rear {
		curValue := decimal.NewFromFloat(q.arr[idx].NewestPrice)
		if max.LessThan(curValue) {
			max = curValue
		}
		if curValue.LessThan(min) {
			min = curValue
		}
		idx = (idx + 1) % q.maxLen
	}
	q.highestPrice = max
	q.lowestPrice = min
}

func (q *QuoteQueue) List() {
	var idx = q.front
	for idx != q.rear {
		log.Println(q.arr[idx])
		idx = (idx + 1) % q.maxLen
	}
}
