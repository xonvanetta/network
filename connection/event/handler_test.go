package event

import (
	"fmt"
	"testing"

	"github.com/xonvanetta/network/connection/utils"

	"github.com/stretchr/testify/assert"
)

func TestHandlerAddAndDo(t *testing.T) {
	counter := utils.NewCounter()
	event := New("uuid", nil)

	tests := []struct {
		name string
		add  func(*testing.T, Handler)
		do   func(*testing.T, Handler)
	}{
		{
			name: "Add",
			add: func(t *testing.T, handler Handler) {
				handler.Add(1, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
			},
			do: func(t *testing.T, handler Handler) {
				err := handler.Do(1, event)
				assert.NoError(t, err)
				counter.Verify(t, 1)
			},
		},
		{
			name: "Multiple adds",
			add: func(t *testing.T, handler Handler) {
				handler.Add(1, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				handler.Add(1, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
			},
			do: func(t *testing.T, handler Handler) {
				err := handler.Do(1, event)
				assert.NoError(t, err)
				counter.Verify(t, 2)
			},
		},
		{
			name: "Multiple adds with error",
			add: func(t *testing.T, handler Handler) {
				handler.Add(1, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				handler.Add(1, func(event Event) error {
					counter.Inc()
					return fmt.Errorf("boom")
				})
			},
			do: func(t *testing.T, handler Handler) {
				err := handler.Do(1, event)
				assert.Error(t, err)
				counter.Verify(t, 2)
			},
		},
		{
			name: "Multiple adds with different types",
			add: func(t *testing.T, handler Handler) {
				handler.Add(1, func(event Event) error {
					counter.Inc()
					assert.Equal(t, "uuid", event.UUID())
					return nil
				})
				handler.Add(2, func(event Event) error {
					counter.Inc()
					return nil
				})
			},
			do: func(t *testing.T, handler Handler) {
				err := handler.Do(1, event)
				assert.NoError(t, err)
				counter.Verify(t, 1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counter.Reset()
			handler := NewHandler()
			test.add(t, handler)
			test.do(t, handler)
		})
	}
}

// 50-60 ns
func BenchmarkAdd(b *testing.B) {
	handler := NewHandler()
	for n := 0; n < b.N; n++ {
		handler.Add(1, func(event Event) error {
			return nil
		})
	}
}

// 670 - 700 ns
func BenchmarkDo(b *testing.B) {
	handler := NewHandler()
	handler.Add(1, func(event Event) error {
		return nil
	})

	event := New("uuid", nil)

	for n := 0; n < b.N; n++ {
		err := handler.Do(1, event)
		assert.NoError(b, err)
	}
}
