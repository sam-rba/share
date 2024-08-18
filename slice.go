package share

// ConstSlice is a read-only slice that can be shared between goroutines. It should be closed after use.
type ConstSlice[T any] struct {
	requests chan<- chan T
}

// NewConstSlice initializes a ConstSlice with the given slice. The slice is not copied, and thus should not be used
// again after creating the ConstSlice.
func NewConstSlice[T any](slc []T) ConstSlice[T] {
	requests := make(chan chan T)
	go serve(slc, requests)
	return ConstSlice[T]{requests}
}

// Elems returns a channel that yields each successive element of the slice.
// Once drained, the channel is closed automatically, and therefore should NOT be closed again by the caller.
func (cs ConstSlice[T]) Elems() <-chan T {
	c := make(chan T) // will be closed by request handler
	cs.requests <- c
	return c
}

func (cs ConstSlice[T]) Close() {
	close(cs.requests)
}

func serve[T any](slc []T, requests <-chan chan T) {
	for request := range requests {
		go handle(request, slc)
	}
}

func handle[T any](request chan<- T, slc []T) {
	defer close(request)
	for _, v := range slc {
		request <- v
	}
}
