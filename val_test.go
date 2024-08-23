package share_test

import (
	"testing"

	"git.samanthony.xyz/share"
)

// Set value in local goroutine, verify in remote goroutine.
func TestValSetLocal(t *testing.T) {
	val := "foo"
	sharedVal := share.NewVal[string]()
	sharedVal.Set <- val
	verifySameVal(sharedVal, val, t)
	go func() {
		defer sharedVal.Close()
		verifySameVal(sharedVal, val, t)
	}()
}

// Set value in remote goroutine, verify in local goroutine.
func TestValSetRemote(t *testing.T) {
	val := "foo"
	sharedVal := share.NewVal[string]()
	defer sharedVal.Close()
	done := make(chan bool)
	defer close(done)
	go func() {
		sharedVal.Set <- val
		verifySameVal(sharedVal, val, t)
		done <- true
	}()
	verifySameVal(sharedVal, val, t)
	<-done
}

func verifySameVal[T comparable](sv share.Val[T], v T, t *testing.T) {
	ret := sv.Get()
	if ret != v {
		t.Errorf("Val.Get() = %v; expected %v", ret, v)
	}
}
