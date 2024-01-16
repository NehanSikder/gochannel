// Contains testing code

package main

import (
	"container/list"
	"fmt"
	"sync"
)

// TODO import channel implementation
// TODO write the semaphore struct
// record how many permits we have left
// use condition variable to help wait when we dont have enough permits
type Semaphore struct {
	permits int
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock()
	for rw.permits <= 0 {
		rw.cond.Wait()
	}
	rw.permits--
	rw.cond.L.Unlock()

}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()
	rw.permits++
	rw.cond.Signal()
	rw.cond.L.Unlock()

}

type Channel[M any] struct {
	capacitySema *Semaphore
	sizeSema     *Semaphore
	mutex        sync.Mutex
	buffer       *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		capacitySema: NewSemaphore(capacity),
		sizeSema:     NewSemaphore(0),
		buffer:       list.New(),
	}
}

// Send function
// Block go routes if the buffer is full
// add message to buffer
// if any receiver go routeine is blocked, resume one of them

func (c *Channel[M]) Send(message M) {
	c.capacitySema.Acquire()
	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()
	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	c.capacitySema.Release()
	c.sizeSema.Acquire()
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()
	return v
}

func findFactors(number int) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

func main() {

	// TODO write test using channel
	// TODO create channel
	// TODO call find factors using channel
	// TODO call find factors regularly
	channel := NewChannel[[]int](0)
	go func() {
		channel.Send(findFactors(4567823))
	}()
	fmt.Println(findFactors(4567823))
	fmt.Println(channel.Receive())
	fmt.Println("Testing completed")

}
