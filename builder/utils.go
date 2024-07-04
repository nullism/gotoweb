package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

func (b *Builder) CheckDirectories() error {
	// Source directory
	sd, err := os.Stat(b.site.SourceDir)
	if err != nil || !sd.IsDir() {
		log.Error("source directory inaccessible", "directory", b.site.SourceDir)
		return err
	}

	// Public directory
	pd, err := os.Stat(b.site.PublicDir)
	if err != nil {
		err2 := os.Mkdir(b.site.PublicDir, 0755)
		if err2 != nil {
			log.Error("could not stat or create public directory", "directory", b.site.PublicDir, "mkdir error", err2, "error", err)
			return err2
		}
	} else if !pd.IsDir() {
		log.Error("public exists and is not a directory", "directory", b.site.PublicDir)
		return fmt.Errorf("public exists and is not a directory")
	}
	return nil
}

// RemovePublic removes cached files from the public directory.
func (b *Builder) RemovePublic() error {
	var err error
	if _, err := os.Stat(b.site.PublicDir); os.IsNotExist(err) {
		return nil
	}

	// remove all children aside from public (for volume mounts)
	err = filepath.WalkDir(b.site.PublicDir, func(path string, info os.DirEntry, err error) error {
		if path == b.site.PublicDir {
			return nil
		}
		err2 := os.RemoveAll(path)
		if err2 != nil {
			log.Error("could not remove public directory", "directory", b.site.PublicDir, "error", err)
			return err2
		}
		return nil
	})
	return err
}
