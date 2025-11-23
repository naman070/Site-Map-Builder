package queue

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}

func (q *Queue[T]) IsEmpty() bool {
	return q.Size() == 0
}

func (q *Queue[T]) Front() (T, bool) {
	var item T
	if q.IsEmpty() {
		return item, false
	}
	return q.items[0], true
}

func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() (T, bool) {
	var item T
	if q.IsEmpty() {
		return item, false
	}
	item = q.items[0]
	q.items = q.items[1:]
	return item, true
}
