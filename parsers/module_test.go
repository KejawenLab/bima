package parsers

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Module_Parse(t *testing.T) {
	cwd, _ := os.Getwd()

	assert.Equal(t, 1, len(ParseModule(fmt.Sprintf("%s/stubs/valids", cwd))))
	assert.Equal(t, 0, len(ParseModule(fmt.Sprintf("%s/stubs/not_exists", cwd))))
	assert.Equal(t, 0, len(ParseModule(fmt.Sprintf("%s/stubs/invalids", cwd))))
}
