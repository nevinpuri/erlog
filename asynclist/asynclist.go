package asynclist

import (
	"erlog/models"
	"sync"
)

// TODO: make v type generic
type AsyncList struct {
	mu		sync.Mutex
	len		int
	size	int
	v		[]models.ErLog
}

// TODO: make two batch sizes, one which the array is created, and one which the array physically can't hold more

func New(size int) AsyncList {
	return AsyncList {v: make([]models.ErLog, size), size: size, len: -1}
}

// Appends a value to the list
func (c *AsyncList) Append(value models.ErLog) {
	c.mu.Lock()

	c.len += 1

	if c.v == nil {
		c.v = make([]models.ErLog, 1)
	}

	if c.len >= c.size {
		c.v = append(c.v, value)
	} else {
		c.v[c.len] = value
	}

	c.mu.Unlock()
}

// Gets a specific value from the list
func (c *AsyncList) Value(idx uint) models.ErLog {
	return c.v[idx]
}

// clears the list
// freeing all memory
func (c *AsyncList) Clear() {
	c.mu.Lock()

	c.v = nil

	c.len = -1
	c.mu.Unlock()
}

// Gets the len of list
func (c *AsyncList) Len() int {
	return c.len
}

// Returns all items in list
// creating a copy first
// no idea why I need to make a copy
func (c *AsyncList) All() []models.ErLog {
	d := make([]models.ErLog, c.len + 1)

	for i := 0; i < c.len + 1; i++ {
		// d[i] = make(models.ErLog, len(c.v[i]))
		d[i] = c.v[i]
		// copy(d[i], c.v[i])
	}

	return d

	// write tests
}
