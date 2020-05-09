// Limit the number of Goroutines. Use algorithm like token bucket.

package glimit

import (
	"reflect"
)

type Glimit struct {
	n int
	c chan struct{}
}

func New(n int) *Glimit {
	return &Glimit{
		n: n,
		c: make(chan struct{}, n),
	}
}

// Run f in a new goroutine but with limit.
func (g *Glimit) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}

// A convenient version of Run
func (g *Glimit) Runf(f interface{}, args ...interface{}) {
	g.c <- struct{}{}
	go func() {
		call(f, args...)
		<-g.c
	}()
}

func call(f interface{}, args ...interface{}) {
	fv := reflect.ValueOf(f)
	in := make([]reflect.Value, len(args))
	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	fv.Call(in)
}
