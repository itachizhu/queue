package queue

import "sync"

type SyncQueue struct {
	lock sync.RWMutex
	cond *sync.Cond
	buffer *Queue
}

func NewSyncQueue() *SyncQueue {
	ch := &SyncQueue{
		buffer: New(),
	}
	ch.cond = sync.NewCond(&ch.lock)
	return ch
}

func (q *SyncQueue) Pop() interface{} {
	c := q.cond
	buffer := q.buffer
	q.lock.RLock()
	defer q.lock.RUnlock()
	for buffer.Length() == 0 {
		c.Wait()
	}
	v := buffer.Peek()
	buffer.Remove()
	return v
}

func (q *SyncQueue) TryPop() (interface{}, bool) {
	buffer := q.buffer
	q.lock.RLock()
	defer q.lock.RUnlock()
	if buffer.Length() > 0 {
		v := buffer.Peek()
		buffer.Remove()
		return v, true
	} else {
		return nil, false
	}
}

func (q *SyncQueue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.buffer.Add(v)
	q.cond.Signal()
}