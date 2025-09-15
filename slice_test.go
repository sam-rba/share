package share_test

import (
	"sync"
	"testing"

	"github.com/sam-rba/share"
)

func TestConstSlice(t *testing.T) {
	orig := []string{"foo", "bar", "baz"}

	shared := share.NewConstSlice(orig)
	defer shared.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		verifySameSlice(shared, orig, t)
		wg.Done()
	}()
	go func() {
		verifySameSlice(shared, orig, t)
		wg.Done()
	}()
	wg.Wait()
}

func verifySameSlice[T comparable](cs share.ConstSlice[T], s []T, t *testing.T) {
	i := 0
	for elem := range cs.Elems() {
		if i < len(s) {
			if elem != s[i] {
				t.Errorf("ConstSlice[%d] = %v; expected %v", i, elem, s[i])
			}
		}
		i++
	}
	if i != len(s) {
		t.Errorf("ConstSlice.Elems() returned %d elements; expected %d", i, len(s))
	}
}
