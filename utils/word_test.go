package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Word_All(t *testing.T) {
	assert.Equal(t, "surya_makan_sabun", Underscore("Surya Makan Sabun"))
	assert.Equal(t, "SuryaMakanSabun", Camelcase("Surya Makan Sabun"))
	assert.Equal(t, "SuryaMakanSabun", Camelcase("surya makan sabun"))
}
