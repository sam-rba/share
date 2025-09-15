package share

// Deque is a double-ended queue with an unlimited capacity.
type Deque[T any] struct {
	// Sending to PutTail adds an element to the back of the deque
	// and never blocks.
	PutTail chan<- T

	// Sending to PutHead adds an element to the front of the deque
	// and never blocks.
	PutHead chan<- T

	// Receiving from TakeHead removes an element from the front
	// of the deque, or, if the queue is empty, blocks until an element
	// is enqueued.
	TakeHead <-chan T

	// Receiving from TakeTail removes an element from the back
	// of the deque, or, if the deque is empty, blocks until an element
	// is enqueued.
	TakeTail <-chan T
}

func NewDeque[T any]() Deque[T] {
	putTail, putHead := make(chan T), make(chan T)
	takeHead, takeTail := make(chan T), make(chan T)

	go run(putTail, putHead, takeHead, takeTail)

	return Deque[T]{
		PutTail:  putTail,
		PutHead:  putHead,
		TakeHead: takeHead,
		TakeTail: takeTail,
	}
}

func run[T any](putTail, putHead <-chan T, takeHead, takeTail chan<- T) {
	defer close(takeTail)
	defer close(takeHead)

	var buf []T

	// While the Put channels are open,
	for ok := true; ok; {
		if len(buf) > 0 {
			buf, ok = putOrTake(putTail, putHead, takeHead, takeTail, buf)
		} else {
			buf, ok = put(putTail, putHead, buf)
		}
	}

	flush(takeHead, takeTail, buf)
}

func flush[T any](takeHead, takeTail chan<- T, buf []T) {
	for len(buf) > 0 {
		select {
		case takeHead <- buf[0]:
			buf = buf[1:]
		case takeTail <- buf[len(buf)-1]:
			buf = buf[:len(buf)-1]
		}
	}
}

func putOrTake[T any](putTail, putHead <-chan T, takeHead, takeTail chan<- T, buf []T) ([]T, bool) {
	select {
	case takeHead <- buf[0]:
		buf = buf[1:]

	case takeTail <- buf[len(buf)-1]:
		buf = buf[:len(buf)-1]

	case v, ok := <-putTail:
		if !ok {
			return buf, false
		}
		buf = append(buf, v)

	case v, ok := <-putHead:
		if !ok {
			return buf, false
		}
		buf = append([]T{v}, buf...)
	}

	return buf, true
}

func put[T any](putTail, putHead <-chan T, buf []T) ([]T, bool) {
	var v T
	var ok bool

	select {
	case v, ok = <-putTail:
		if ok {
			buf = append(buf, v)
		}
	case v, ok = <-putHead:
		if ok {
			buf = append([]T{v}, buf...)
		}
	}

	return buf, ok
}

// Close the Put channels of the deque.
// The deque will wait until all elements have been drained through
// either of the Take channels before closing them.
func (dq Deque[T]) Close() {
	close(dq.PutTail)
	close(dq.PutHead)
}
