package fsys

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

// FileSystem is an interface that abstracts file system operations.
type FileSystem interface {
	Abs(string) (string, error)
	Base(string) string
	Copy(string, string, os.FileMode) error
	Dir(string) string
	Exists(string) bool
	Ext(string) string
	FindInParent(string, int) (string, error)
	Join(elem ...string) string
	IsAbs(string) bool
	Mkdir(string, os.FileMode) error
	MkdirAll(string, os.FileMode) error
	ReadFile(string) ([]byte, error)
	RemoveAll(string) error
	Stat(string) (os.FileInfo, error)
	WalkDir(string, fs.WalkDirFunc) error
	WriteFile(string, []byte, os.FileMode) error
}

// OsFileSystem implements the FileSystem interface using the os package.
type OsFileSystem struct {
}

// Abs returns an absolute representation of path.
func (OsFileSystem) Abs(name string) (string, error) {
	return filepath.Abs(name)
}

// Base returns the last element of path (basename).
func (OsFileSystem) Base(name string) string {
	return filepath.Base(name)
}

// Copy copies recursively a path from src to dst with the given permissions.
func (OsFileSystem) Copy(src, dst string, perm os.FileMode) error {
	return cp.Copy(src, dst, cp.Options{AddPermission: perm})
}

// Dir returns all but the last element of path, typically the path's directory.
func (OsFileSystem) Dir(name string) string {
	return filepath.Dir(name)
}

// Exists returns true if the path exists.
func (OsFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Ext returns the file name extension used by path.
func (OsFileSystem) Ext(name string) string {
	return filepath.Ext(name)
}

// FindInParent searches for a file in the current directory and then in [maxDirCount] parent directories.
func (o *OsFileSystem) FindInParent(name string, maxDirCount int) (string, error) {
	relative := "."
	for range maxDirCount {
		relative = o.Join(relative, "..")
	}

	p := o.Join(relative, name)
	_, err := o.Stat(p)
	if err != nil {
		if maxDirCount < 1 {
			return "", fmt.Errorf("could not find %v in self or parent directories: %w", name, err)
		}
		return o.FindInParent(name, maxDirCount-1)
	}
	fp, err := o.Abs(p)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path for %v: %w", name, err)
	}
	return fp, nil
}

// IsAbs reports whether the path is absolute.
func (OsFileSystem) IsAbs(name string) bool {
	return filepath.IsAbs(name)
}

// Join joins any number of path elements into a single path, adding an os-specific separator if necessary.
func (OsFileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// MkDir creates a new directory with the specified name and permission bits.
func (OsFileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
func (OsFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// ReadFile reads the file named by filename and returns the contents.
func (OsFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// RemoveAll removes all files and directories in the named directory, including the root.
func (OsFileSystem) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

// Stat returns an os.FileInfo describing the named file or error.
func (OsFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root.
func (OsFileSystem) WalkDir(name string, walkFn fs.WalkDirFunc) error {
	return filepath.WalkDir(name, walkFn)
}

// WriteFile writes data to a file named by filename with [perm].
func (OsFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
