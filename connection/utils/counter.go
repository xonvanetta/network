package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ONLY USED IN TESTING
type Counter interface {
	Inc()
	Verify(t *testing.T, amount int)
	Reset()
}

type counter struct {
	mutex  *sync.Mutex
	amount int
}

func NewCounter() Counter {
	return &counter{
		mutex: &sync.Mutex{},
	}
}

func (c *counter) Inc() {
	c.mutex.Lock()
	c.amount++
	c.mutex.Unlock()
}

func (c *counter) Verify(t *testing.T, amount int) {
	c.mutex.Lock()
	assert.Equal(t, amount, c.amount, "amount of calls don't match with the counter")
	c.mutex.Unlock()
}

func (c *counter) Reset() {
	c.mutex.Lock()
	c.amount = 0
	c.mutex.Unlock()
}
