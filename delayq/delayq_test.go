package delayq

import (
	"testing"
	"time"
)

func TestQ(t *testing.T) {
	myq := New()
	now := time.Now()
	myq.Put(now, 0)
	myq.Put(now.Add(time.Second), 1)
	myq.Put(now.Add(3*time.Second), 4)

	go func() {
		time.Sleep(time.Second * 2)
		myq.Put(now, 2)
		myq.Put(now, 3)
	}()

	for i := 0; i < 5; i++ {
		x := myq.Get()
		if i != x.(int) {
			t.Error("i != x", i)
		}
	}
}
