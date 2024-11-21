package share_test

import (
	"testing"

	"github.com/sam-rba/share"
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

// Val.TryGet() before Set should fail.
func TestValTryGetFail(t *testing.T) {
	sv := share.NewVal[int]() // type is arbitrary
	defer sv.Close()
	if v, ok := sv.TryGet(); ok {
		t.Errorf("Val.TryGet() succeeded (returned %v) before value was set; expected to fail", v)
	}
}

// Val.TryGet() after Set should succeed.
func TestValTryGet(t *testing.T) {
	sv := share.NewVal[string]()
	defer sv.Close()
	v := "foo"
	sv.Set <- v
	ret, ok := sv.TryGet()
	if !ok {
		t.Error("Val.TryGet() failed")
	}
	if *ret != v {
		t.Errorf("Val.TryGet() returned %v; expected %v", ret, v)
	}
}

func verifySameVal[T comparable](sv share.Val[T], v T, t *testing.T) {
	ret := sv.Get()
	if ret != v {
		t.Errorf("Val.Get() = %v; expected %v", ret, v)
	}
}
