package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser(t *testing.T) {
	gen := UserGenID{lastID: 0}
	for i := 1; i <= 10; i++ {
		id := gen.GetID()
		assert.Equal(t, i, id)
	}
}
