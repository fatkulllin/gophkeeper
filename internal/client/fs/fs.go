package fs

import (
	"os"
	"path/filepath"
)

func PrepareAppDir() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(baseDir, "gophkeeper")

	if err := os.MkdirAll(path, 0700); err != nil {
		return "", err
	}

	return path, nil
}
