package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Type_All(t *testing.T) {
	types := NewType()

	assert.Equal(t, 15, len(types.List()))
	assert.Equal(t, "float64", types.Value("double"))
	assert.Equal(t, "", types.Value("not_exist"))
}
