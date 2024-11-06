package commands

import (
	"log"
)

type Queue struct {
	list []string
}

func NewQueue() {

}
func (q *Queue) Push(s string) {
	q.list = append(q.list, s)

	log.Printf("Push: Length of list is now %d", len(q.list))
}

func (q *Queue) Pop() string {
	x := q.list[0]

	q.list = q.list[1:]
	log.Printf("Pop: Length of list is now %d", len(q.list))
	return x
}

func (q *Queue) IsEmpty() bool {
	log.Printf("Checking for emptiness: Length of list is now %d", len(q.list))
	if len(q.list) == 0 {
		return true
	} else {
		return false
	}
}
