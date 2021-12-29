package ciperutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	data := "data"
	assert.NotEqual(t, "data", Hash(data))
}

func TestCheckHash(t *testing.T) {
	exp := "a17c9aaa61e80a1bf71d0d850af4e5baa9800bbd"

	assert.True(t, CheckHash(exp, "data"))
}
