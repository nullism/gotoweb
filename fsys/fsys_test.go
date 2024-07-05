package fsys

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindInParent(t *testing.T) {
	fs := &OsFileSystem{}

	path, err := fs.FindInParent("files_that_does_not_exist", 3)
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Empty(t, path)

	path, err = fs.FindInParent("go.mod", 5)
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
}

func Test_filepathFunctions(t *testing.T) {
	fs := &OsFileSystem{}

	assert.Equal(t, "/foo/bar", fs.Join("/foo", "bar"))
	assert.True(t, fs.IsAbs("/foo/bar"))
	assert.False(t, fs.IsAbs("foo/bar"))
	assert.Equal(t, "bar.txt", fs.Base("/foo/bar.txt"))
	assert.Equal(t, "/foo", fs.Dir("/foo/bar"))
	assert.Equal(t, ".txt", fs.Ext("foo/bar.baz.txt"))
}
