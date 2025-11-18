package filemanager

import (
	"os"
	"path/filepath"
)

type FileManager struct {
}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f *FileManager) SaveFile(filename string, body string, permission os.FileMode) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "gophkeeper")
	err = os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}
	filePath := filepath.Join(path, filename)
	err = os.WriteFile(filePath, []byte(body), permission)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileManager) LoadFile(filename string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "gophkeeper", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (f *FileManager) RemoveFile(filename string) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "gophkeeper", filename)
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
