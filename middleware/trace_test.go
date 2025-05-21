package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceMiddleware_generateTraceID(t *testing.T) {
	var middleware TraceMiddleware
	t.Run("UUID uniqueness", func(t *testing.T) {
		generated := make(map[string]struct{})
		loopCount := 1000
		for i := 0; i < loopCount; i++ {
			got := middleware.generateTraceID()
			assert.Equal(t, len(got), 36)
			if _, exists := generated[got]; exists {
				t.Errorf("Duplicate UUID detected: %v", got)
				return
			}
			generated[got] = struct{}{}
		}
		assert.Equal(t, loopCount, len(generated))
	})
}
