package commands

import (
	"log"
	"sync"
)

type Queue struct {
	mu   sync.Mutex
	list []string
}

func NewQueue() {

}
func (q *Queue) Push(s string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list = append(q.list, s)
	log.Printf("Push: Length of list is now %d", len(q.list))
}

func (q *Queue) Pop() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	x := q.list[0]
	q.list = q.list[1:]

	log.Printf("Pop: Length of list is now %d", len(q.list))
	return x
}

func (q *Queue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	log.Printf("Checking for emptiness: Length of list is now %d", len(q.list))
	if len(q.list) == 0 {
		return true
	} else {
		return false
	}
}

func (q *Queue) PrintQueue() {
	for i, elem := range q.list {
		log.Printf("Elem Number: %d, Song: %s", i, elem)
	}
}
