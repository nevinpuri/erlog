package asynclist

import (
	"sync"
)

// TODO: make v type generic
type AsyncList struct {
	mu		sync.Mutex
	len		int
	size	int
	v		[][]byte
}

// TODO: make two batch sizes, one which the array is created, and one which the array physically can't hold more

func New(size int) AsyncList {
	return AsyncList {v: make([][]byte, size), size: size, len: -1}
}

// Appends a value to the list
func (c *AsyncList) Append(value []byte) {
	c.mu.Lock()

	c.len += 1

	if c.len >= c.size {
		c.v = append(c.v, value)
	} else {
		c.v[c.len] = value
	}

	c.mu.Unlock()
}

// Gets a specific value from the list
func (c *AsyncList) Value(idx uint) []byte {
	return c.v[idx]
}

// clears the list
// freeing all memory
func (c *AsyncList) Clear() {
	c.mu.Lock()

	for i := 0; i < c.len + 1; i++ {
		c.v[i] = nil
	}

	c.len = -1
	c.mu.Unlock()
}

// Gets the len of list
func (c *AsyncList) Len() int {
	return c.len
}

// Returns all items in list
// creating a copy first
func (c *AsyncList) All() [][]byte {
	d := make([][]byte, c.len + 1)

	for i := 0; i < c.len + 1; i++ {
		d[i] = make([]byte, len(c.v[i]) + 1)
		copy(d[i], c.v[i])
	}

	return d
}
