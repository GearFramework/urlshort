package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	gen := UserGenID{lastID: 0}
	for i := 1; i <= 10; i++ {
		id := gen.GetID()
		assert.Equal(t, i, id)
	}
}
