package parsers

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Route_Parse(t *testing.T) {
	parser := Route{}

	cwd, _ := os.Getwd()

	assert.Equal(t, 2, len(parser.Parse(fmt.Sprintf("%s/stubs/valids", cwd))))
	assert.Equal(t, 0, len(parser.Parse(fmt.Sprintf("%s/stubs/not_exists", cwd))))
	assert.Equal(t, 0, len(parser.Parse(fmt.Sprintf("%s/stubs/invalids", cwd))))
}
