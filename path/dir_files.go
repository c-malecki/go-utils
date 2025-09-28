package path

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// Will get all file paths for the specified file extension in the specified target directory
func PathsForFilesInDir(dir, ext string) ([]string, error) {
	var paths []string
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ext) {
			return nil
		}

		paths = append(paths, path)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("filepath.WalkDir %w", err)
	}

	return paths, nil
}
