package trylock

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestMutexLock(t *testing.T) {
	mu := New()

	mu.Lock()
	mu.Unlock()
	mu.Lock()
	mu.Unlock()

	mu.RLock()
	mu.RUnlock()
	mu.RLock()
	mu.RUnlock()

	mu.TryLock(nil)
	mu.Unlock()
	mu.TryLockTimeout(5 * time.Second)
	mu.Unlock()

	mu.RTryLockTimeout(0)
	mu.RUnlock()
	mu.RTryLockTimeout(5 * time.Second)
	mu.RUnlock()
}

func TestMutexLockTryLock(t *testing.T) {
	mu := New()

	if ok := mu.TryLock(nil); !ok {
		t.Errorf("cannot Lock !!!")
	}
	if ok := mu.TryLock(nil); ok {
		t.Errorf("cannot Lock twice !!!")
	}

	mu.Unlock()
}

func TestMutexLockAfterUnlock(t *testing.T) {
	mu := New()
	mu.Lock()

	go func() {
		time.Sleep(50 * time.Millisecond)
		mu.Unlock()
	}()

	mu.Lock()
	mu.Unlock()
}

func TestMutexLockAfterRUnlock(t *testing.T) {
	mu := New()
	mu.RLock()

	go func() {
		time.Sleep(50 * time.Millisecond)
		mu.RUnlock()
	}()

	mu.Lock()
	mu.Unlock()
}

func TestMutexLockTryLockTimeout(t *testing.T) {
	mu := New()
	mu.Lock()

	if ok := mu.TryLockTimeout(10 * time.Millisecond); ok {
		t.Errorf("should not Lock in 10ms !!!")
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		mu.Unlock()
	}()
	if ok := mu.TryLockTimeout(200 * time.Millisecond); !ok {
		t.Errorf("cannot Lock after 200ms !!!")
	}

	mu.Unlock()
}

func TestMutexLockRTryLockTimeout(t *testing.T) {
	mu := New()
	mu.Lock()

	if ok := mu.RTryLockTimeout(10 * time.Millisecond); ok {
		t.Errorf("should not Lock in 10ms !!!")
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		mu.Unlock()
	}()
	if ok := mu.RTryLockTimeout(200 * time.Millisecond); !ok {
		t.Errorf("cannot Lock after 200ms !!!")
	}
	mu.RUnlock()
}

func TestMutexLockUnLockTwice(t *testing.T) {
	mu := New()
	mu.Lock()
	defer func() {
		if x := recover(); x != nil {
			if x != "Unlock() failed" {
				t.Errorf("unexpect panic")
			}
		} else {
			t.Errorf("should panic after unlock twice")
		}
	}()
	mu.Unlock()
	mu.Unlock()
}

func TestMutexLockRLockTwice(t *testing.T) {
	mu := New()
	mu.RLock()
	mu.RLock()
	mu.RUnlock()
	mu.RUnlock()
}

func TestMutexLockUnLockInvalid(t *testing.T) {
	mu := New()
	mu.Lock()
	defer func() {
		if x := recover(); x != nil {
			if x != "RUnlock() failed" {
				t.Errorf("unexpect panic")
			}
		} else {
			t.Errorf("should panic after RUnlock a write lock")
		}
	}()
	mu.RUnlock()
}

func TestMutexLockBroadcast(t *testing.T) {
	mu := New()
	mu.Lock()

	done := int32(0)
	for i := 0; i < 3; i++ {
		go func() {
			mu.RLock()
			atomic.AddInt32(&done, 1)
			mu.RUnlock()
		}()
	}

	time.Sleep(10 * time.Millisecond)

	mu.Unlock()

	time.Sleep(10 * time.Millisecond)

	if atomic.LoadInt32(&done) != 3 {
		t.Fatal("Broadcast is failed")
	}
}
