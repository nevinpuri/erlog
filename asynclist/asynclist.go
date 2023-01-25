package asynclist

import "sync"

// TODO: make v type generic
type AsyncList struct {
	mu		sync.Mutex
	len		int
	size	int
	v		[][]byte
}

func New(size int) AsyncList {
	return AsyncList {v: make([][]byte, size), len: 0}
}

// Appends a value to the list
func (c *AsyncList) Append(value []byte) {
	c.mu.Lock()
	c.v = append(c.v, value)
	c.len += 1
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
	c.v = nil
	c.len = 0
	c.mu.Unlock()
}

// Gets the len of list
func (c *AsyncList) Len() int {
	return c.len
}

// Returns all items in list
// creating a copy first
func (c *AsyncList) All() [][]byte {
	d := make([][]byte, len(c.v))

	for i := range c.v {
		d[i] = make([]byte, len(c.v[i]))
		copy(d[i], c.v[i])
	}

	return d
}
