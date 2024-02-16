package add

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Run("when the db has a nil value, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Add(nil)
		})
	})

	t.Run("when the task name is missing we should return a 400 error", func(t *testing.T) {
		Add(nil)
	})
}
