package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTestFile(t *testing.T) {
	assert := require.New(t)

	fileName, rmFn := NewTestFile(t, 0600, "line1\n", "line2\n")
	assert.NotEmpty(fileName)

	_, err := os.Stat(fileName)
	assert.NoError(err)

	content, err := ioutil.ReadFile(fileName)
	assert.NoError(err)
	assert.Contains(string(content), "line2")

	rmFn()

	_, err = os.Stat(fileName)
	assert.EqualError(err, fmt.Sprintf("stat %s: no such file or directory", fileName))
}
