package glimit

import (
	"testing"
	"time"
)

func TestGlimit(t *testing.T) {
	gl := New(3)
	c := [4]bool{}
	block := make(chan struct{})
	go func() {
		for i := 0; i < 4; i++ {
			j := i
			f := func() {
				c[j] = true
				block <- struct{}{}
			}
			gl.Run(f)
		}
	}()
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 3; i++ {
		if !c[i] {
			t.Errorf("c[%d] should be true", i)
		}
	}
	if c[3] {
		t.Errorf("c[3] should be false")
	}
	for i := 0; i < 3; i++ {
		<-block
	}
	time.Sleep(10 * time.Millisecond)
	if !c[3] {
		t.Errorf("c[3] should be true")
	}
}
