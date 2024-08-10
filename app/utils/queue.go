package utils

// Node represents a single element in the queue
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}

// Queue is a data structure that stores items and allows us access in a First-In/First-Out manner.
type Queue[T any] struct {
	Head *Node[T]
	Tail *Node[T]
}

// Add an item to the end of the queue
func (q *Queue[T]) Enqueue(value T) {
	newNode := NewNode(value)

	if q.Tail != nil {
		q.Tail.Next = newNode
	}
	q.Tail = newNode

	if q.Head == nil {
		q.Head = newNode
	}
}

// Remove the first (oldest) item
func (q *Queue[T]) Dequeue() *Node[T] {

	if q.Head == nil {
		return nil
	}

	oldHead := q.Head
	q.Head = q.Head.Next

	if q.Head == nil {
		q.Tail = nil
	}

	return oldHead
}

// View the first (oldest) item
func (q *Queue[T]) Peek() *Node[T] {
	return q.Head
}

// Delete Queue
func (q *Queue[T]) Delete() {
	q.Head = nil
	q.Tail = nil

}

// Convert a list of items into a queue
func (q *Queue[T]) ArrayToQueue(values []T) *Queue[T] {
	for _, value := range values {
		q.Enqueue(value)
	}
	return q

}

// Initializes a new queue instance and returns a pointer to it.
func NewQueue[T any](values []T) *Queue[T] {
	q := &Queue[T]{}
	return q.ArrayToQueue(values)
}
