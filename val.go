package share

// Val is a concurrent interface to a piece of shared data.
//
// A client can read the data by sending a channel via request, and the stored value will
// be sent back via the channel. The client is responsible for closing the channel.
//
// The stored value can be changed by sending the new value via Set. Requests block until
// the first value is received on Set.
//
// Val should be closed after use.
type Val[T any] struct {
	Request chan<- chan T
	Set     chan<- T
}

func NewVal[T any]() Val[T] {
	request := make(chan chan T)
	set := make(chan T)
	go func() {
		val := <-set // wait for initial value
		for {
			select {
			case v, ok := <-set:
				if !ok { // closed
					return
				}
				val = v
			case req, ok := <-request:
				if !ok { // closed
					return
				}
				go func() { // don't wait for client to receive
					req <- val
				}()
			}
		}
	}()
	return Val[T]{request, set}
}

// Get makes a synchronous request and returns the stored value.
func (v Val[T]) Get() T {
	c := make(chan T)
	defer close(c)
	v.Request <- c
	return <-c
}

// TryGet returns the stored value if it has already been set, or false if it hasn't.
func (v Val[T]) TryGet() (*T, bool) {
	c := make(chan T)
	defer close(c)
	select {
	case v.Request <- c:
		val := <-c
		return &val, true
	default:
		return nil, false
	}
}

func (v Val[T]) Close() {
	close(v.Request)
	close(v.Set)
}
