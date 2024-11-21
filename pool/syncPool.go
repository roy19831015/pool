package pool

import (
	"errors"
	"time"
)

type SyncPool[T any] struct {
	chMutex chan T
	factory Factory[T]
}

func (c *SyncPool[T]) Init(cap int, fac Factory[T]) error {
	c.factory = fac
	c.chMutex = make(chan T, cap)
	for i := 0; i < cap; i++ {
		obj, err := c.factory()
		if err != nil {
			return err
		}
		c.chMutex <- obj
	}
	return nil
}

// Get function: when time out reached, Throw out error
func (c *SyncPool[T]) Get(timeout time.Duration) (T, error) {
	var ret T
	select {
	case ret = <-c.chMutex:
		return ret, nil
	case <-time.After(timeout):
		return ret, errors.New("time out reached")
	}
}

func (c *SyncPool[T]) Back(obj T) {
	c.chMutex <- obj
}
