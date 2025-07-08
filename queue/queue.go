package queue

// making the global queue to acess it everywhere
var Que = &Queue{}

type Queue struct {
	items []string
}

func New() *Queue {
	return &Queue{items: []string{}}
}

func (q *Queue) Enqueue(item string) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Len() int {
	return len(q.items)
}
