package share

// Queue is a FIFO queue with an unlimited capacity.
//
// Closing the Enqueue channel closes the Queue. The Queue waits
// until all elements have been drained from the Dequeue channel
// before closing it.
type Queue[T any] struct {
	// Sending to Enqueue adds an element to the back of the Queue
	// and never blocks.
	Enqueue chan<-T


	// Receiving from Dequeue removes an element from the front
	// of the queue or, if the queue is empty, blocks until an element
	// is enqueued.
	Dequeue <-chan T
}

func NewQueue[T any]() Queue[T] {
	in, out := make(chan T), make(chan T)

	go func() {
		defer close(out)

		var queue []T

		for v := range in {
			queue = append(queue, v)

			for len(queue) > 0 {
				select {
				case out <- queue[0]:
					queue = queue[1:]
				case v, ok := <-in:
					if !ok {
						for _, x := range queue {
							out <- x
						}
						return
					}
					queue = append(queue, v)
				}
			}
		}
	}()

	return Queue[T]{Enqueue: in, Dequeue: out}
}
