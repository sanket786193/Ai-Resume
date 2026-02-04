package ocr

import (
	"io"
	"os"
	"path/filepath"
)

// Storage handles file storage operations
type Storage interface {
	SaveFile(src io.Reader, filename string) (string, error)
	DeleteFile(filepath string) error
	FileExists(filepath string) bool
	GetImagePath(filename string) string
}

// fileStorage implements the Storage interface
type fileStorage struct {
	uploadDir string
}

// NewStorage creates a new file storage instance
func NewStorage(uploadDir string) Storage {
	return &fileStorage{
		uploadDir: uploadDir,
	}
}

// SaveFile saves a file to the upload directory
func (s *fileStorage) SaveFile(src io.Reader, filename string) (string, error) {
	// Ensure upload directory exists
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(s.uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = dst.ReadFrom(src); err != nil {
		return "", err
	}

	return filePath, nil
}

// DeleteFile removes a file from the filesystem
func (s *fileStorage) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

// FileExists checks if a file exists
func (s *fileStorage) FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// GetImagePath returns the full path for an image file
func (s *fileStorage) GetImagePath(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}
