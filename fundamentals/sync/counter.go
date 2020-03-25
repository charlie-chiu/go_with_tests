package counter

import "sync"

type Counter struct {
	value int
	mutex sync.Mutex
}

func (c *Counter) Inc() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
