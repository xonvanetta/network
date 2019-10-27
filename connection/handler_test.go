package connection

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xonvanetta/network"
)

type counter struct {
	calls int
}

func (c *counter) Inc() {
	c.calls++
}

func (c *counter) Verify(t *testing.T, calls int) {
	assert.Equal(t, calls, c.calls, "amount of callbacks don't match with the counter")
}

func (c *counter) Reset() {
	c.calls = 0
}

func TestHandlerAddAndDo(t *testing.T) {
	counter := &counter{}

	tests := []struct {
		name string
		add  func(*testing.T)
		do   func(*testing.T)
	}{
		{
			name: "Add",
			add: func(t *testing.T) {
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
			},
			do: func(t *testing.T) {
				err := do(network.Connecting, "uuid", nil)
				assert.NoError(t, err)
				counter.Verify(t, 1)
			},
		},
		{
			name: "Multiple adds",
			add: func(t *testing.T) {
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
			},
			do: func(t *testing.T) {
				err := do(network.Connecting, "uuid", nil)
				assert.NoError(t, err)
				counter.Verify(t, 2)
			},
		},
		{
			name: "Multiple adds with error",
			add: func(t *testing.T) {
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					return fmt.Errorf("boom")
				})
			},
			do: func(t *testing.T) {
				err := do(network.Connecting, "uuid", nil)
				assert.Error(t, err)
				counter.Verify(t, 2)
			},
		},
		{
			name: "Multiple adds with different types",
			add: func(t *testing.T) {
				Add(network.Connecting, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				Add(network.Ping, func(event Event) error {
					counter.Inc()
					return nil
				})
			},
			do: func(t *testing.T) {
				err := do(network.Connecting, "uuid", nil)
				assert.NoError(t, err)
				counter.Verify(t, 1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter.Reset()
			handlers = make(callbacks)
			test.add(t)
			test.do(t)
		})
	}
}

// 50 - 60 ns
func BenchmarkAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Add(network.Ping, func(event Event) error {
			return nil
		})
	}
}

// 700 - 730 ns
func BenchmarkDo(b *testing.B) {
	Add(network.Ping, func(event Event) error {
		return nil
	})

	for n := 0; n < b.N; n++ {
		err := do(network.Ping, "uuid", nil)
		assert.NoError(b, err)
	}
}
