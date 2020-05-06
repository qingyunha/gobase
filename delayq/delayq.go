package delayq

// A delay queue base on min-heap

import (
	"container/heap"
	"sync"
	"time"
)

type item struct {
	t    time.Time
	data interface{}
}

type itemHeap []item

func (h itemHeap) Len() int           { return len(h) }
func (h itemHeap) Less(i, j int) bool { return h[i].t.Before(h[j].t) }
func (h itemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *itemHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}

func (h *itemHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type DelayQ struct {
	h  itemHeap
	mu sync.Mutex

	// notify new item added.
	// condition variable with wait imeout should be more neat.
	// but see https://github.com/golang/go/issues/9578#issuecomment-97145802
	// ensure cap(event) == 1
	// would cause a void notify in some case, but would't loss notify.
	event chan struct{}
}

func New() *DelayQ {
	return &DelayQ{event: make(chan struct{}, 1)}
}

func (q *DelayQ) Put(t time.Time, v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	i := item{t, v}
	heap.Push(&q.h, i)
	select {
	case q.event <- struct{}{}:
	default:
	}
}

func (q *DelayQ) Get() interface{} {
	for {
		q.mu.Lock()
		if len(q.h) == 0 {
			q.mu.Unlock()
			<-q.event
			continue
		}

		now := time.Now()
		wait := q.h[0].t.Sub(now)
		if wait <= 0 {
			x := heap.Pop(&q.h)
			q.mu.Unlock()
			return x.(item).data
		}
		q.mu.Unlock()
		t := time.NewTimer(wait)
		select {
		case <-t.C:
		case <-q.event:
			if !t.Stop() {
				<-t.C
			}
		}
	}
}
