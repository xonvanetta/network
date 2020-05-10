package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool().(*pool)

	pool.Add("f09df0aa-8360-4b96-8041-bb93026ac8a0", nil)
	pool.Add("9e6f7105-9f21-4d16-ba56-a3aee8180163", nil)

	assert.Contains(t, pool.connections, "f09df0aa-8360-4b96-8041-bb93026ac8a0")
	assert.Contains(t, pool.connections, "9e6f7105-9f21-4d16-ba56-a3aee8180163")
	assert.Len(t, pool.connections, 2)

	conns := pool.All()
	assert.Len(t, conns, 2)

	pool.Remove("9e6f7105-9f21-4d16-ba56-a3aee8180163")

	assert.Len(t, conns, 2)
	assert.Len(t, pool.connections, 1)
	conn := pool.Get("f09df0aa-8360-4b96-8041-bb93026ac8a0")
	assert.NotNil(t, conn)
}
