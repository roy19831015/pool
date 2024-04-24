package pool

import (
	"time"
)

type CommonPool[T any] struct {
	chMutex chan T
	factory Factory[T]
}

type Factory[T any] func() (T, error)

func (c *CommonPool[T]) Init(cap int, fac Factory[T]) error {
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

// Get function: when time out reached, Recreate a new T
func (c *CommonPool[T]) Get(timeout time.Duration) (T, error) {
	var ret T
	select {
	case ret = <-c.chMutex:
		return ret, nil
	case <-time.After(timeout):
		ret, err := c.factory()
		if err != nil {
			return ret, err
		}
		return ret, nil
	}
}

func (c *CommonPool[T]) Back(obj T) {
	c.chMutex <- obj
}
