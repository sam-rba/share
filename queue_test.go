package share_test

import (
	"testing"

	"github.com/sam-rba/share"
)

func TestQueue(t *testing.T) {
	q := share.NewQueue[string]()

	vals := []string{"foo", "bar", "baz", "xyz"}
	
	go func() {
		for _, v := range vals {
			q.Enqueue <- v
		}
		close(q.Enqueue)
	}()

	go func() {
		i := 0
		for front := range q.Dequeue {
			t.Log("received", front, "from queue")
			if i > len(vals)-1 {
				t.Fatal("received too many elements from queue")
			}
			if front != vals[i] {
				t.Fatalf("received %v from queue; wanted %v", front, vals[i])
			}
			i++
		}
		if i < len(vals) {
			t.Fatal("did not receive enough values from queue")
		}
	}()
}