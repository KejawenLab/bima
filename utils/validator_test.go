package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	ID   string `validate:"required"`
	Name string `validate:"required"`
}

func Test_Validator(t *testing.T) {
	data1 := Data{
		ID: "test",
	}

	msg, err := Validate(&data1)

	assert.NotNil(t, err)
	assert.NotEmpty(t, msg)

	data2 := Data{
		ID:   "test",
		Name: "test",
	}

	msg, err = Validate(&data2)

	assert.Nil(t, err)
	assert.Empty(t, msg)
}
