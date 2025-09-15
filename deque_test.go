package share_test

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/sam-rba/share"
)

func TestDequeFIFO(t *testing.T) {
	dq := share.NewDeque[string]()
	vals := []string{"foo", "bar", "baz", "xyz"}
	testFIFO(t, dq.PutTail, dq.TakeHead, vals, dq.Close)
}

func TestDequeReverseFIFO(t *testing.T) {
	dq := share.NewDeque[string]()
	vals := []string{"foo", "bar", "baz", "xyz"}
	testFIFO(t, dq.PutHead, dq.TakeTail, vals, dq.Close)
}

func testFIFO[T comparable](t *testing.T, put chan<- T, take <-chan T, vals []T, end func()) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		produce(put, vals)
		end()
		wg.Done()
	}()
	go func() {
		consume(t, take, vals)
		wg.Done()
	}()
	wg.Wait()
}

func TestDequeLIFO(t *testing.T) {
	dq := share.NewDeque[string]()
	vals := []string{"foo", "bar", "baz", "xyz"}
	testLIFO(t, dq.PutTail, dq.TakeTail, vals, dq.Close)
}

func TestDequeReverseLIFO(t *testing.T) {
	dq := share.NewDeque[string]()
	vals := []string{"foo", "bar", "baz", "xyz"}
	testLIFO(t, dq.PutHead, dq.TakeHead, vals, dq.Close)
}

func testLIFO[T comparable](t *testing.T, put chan<- T, take <-chan T, vals []T, end func()) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		produce(put, vals)
		end()
		wg.Done()
	}()

	want := slices.Clone(vals)
	slices.Reverse(want)
	wg.Wait()
	consume(t, take, want)
}

func produce[T any](put chan<- T, vals []T) {
	for _, v := range vals {
		put <- v
	}
}

func consume[T comparable](t *testing.T, take <-chan T, want []T) {
	i := 0
	for v := range take {
		if i >= len(want) {
			t.Fatal("received too many")
		}
		if v != want[i] {
			t.Fatalf("Index %d: %v; want %v", i, v, want[i])
		}
		i++
	}
	if i < len(want) {
		t.Fatalf("Only received %d; want %d", i, len(want))
	}
}

func TestDequePutback(t *testing.T) {
	dq := share.NewDeque[string]()

	dq.PutTail <- "foo"
	dq.PutTail <- "bar"
	dq.PutTail <- "baz"
	dq.PutTail <- "xyz"

	<-dq.TakeHead // foo
	<-dq.TakeHead // bar
	<-dq.TakeHead // baz

	dq.PutHead <- "baz"
	dq.PutHead <- "bar"

	dq.Close()

	consume(t, dq.TakeHead, []string{"bar", "baz", "xyz"})
}

func ExampleDeque() {
	// Use as a FIFO queue
	dq := share.NewDeque[string]()

	// Producer
	go func() {
		defer dq.Close()
		for _, word := range []string{"foo", "bar", "baz"} {
			dq.PutTail <- word
		}
	}()

	for word := range dq.TakeHead {
		fmt.Println(word)
	}
	// Output:
	// foo
	// bar
	// baz
}
