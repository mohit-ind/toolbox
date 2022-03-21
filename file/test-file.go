package file

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func NewTestFile(t *testing.T, perm fs.FileMode, contents ...string) (fileName string, removeFn func()) {
	assert := require.New(t)
	tmpfile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("go-test-%d", time.Now().Unix()))
	assert.NoError(err, "Temp files should shave been created")

	f, err := os.OpenFile(tmpfile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	assert.NoError(err, "Temp file should have been opened for write")

	for _, content := range contents {
		_, err = f.WriteString(content)
		assert.NoError(err, "Content should have been written to the temp file")
	}

	assert.NoError(tmpfile.Close(), "Temp file should have been closed")

	return tmpfile.Name(), func() {
		assert.NoError(os.Remove(tmpfile.Name()))
	}
}
