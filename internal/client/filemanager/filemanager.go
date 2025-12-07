package filemanager

import (
	"os"
	"path/filepath"
)

type FileManager struct {
	cfgDir string
}

func NewFileManager(cfgDir string) *FileManager {
	return &FileManager{cfgDir: cfgDir}
}

func (f *FileManager) SaveFile(filename string, body string, permission os.FileMode) error {
	filePath := filepath.Join(f.cfgDir, filename)
	err := os.WriteFile(filePath, []byte(body), permission)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileManager) LoadFile(filename string) (string, error) {
	path := filepath.Join(f.cfgDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (f *FileManager) RemoveFile(filename string) error {
	path := filepath.Join(f.cfgDir, filename)
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
