package smokepod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindTestFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("accessing path: %w", err)
	}

	if !info.IsDir() {
		if strings.HasSuffix(path, ".test") {
			return []string{path}, nil
		}
		return nil, fmt.Errorf("not a .test file: %s", path)
	}

	var files []string
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(filePath, ".test") {
			files = append(files, filePath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return files, nil
}
