package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage interface {
	Store(ctx context.Context, file io.Reader, filename string) (string, error)
	GetURL(path string) string
	Delete(ctx context.Context, path string) error
}

type LocalFileStorage struct {
	basePath string
	baseURL  string
}

func NewLocalFileStorage(basePath, baseURL string) (*LocalFileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalFileStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

func (s *LocalFileStorage) Store(ctx context.Context, file io.Reader, filename string) (string, error) {
	timestamp := time.Now().Unix()
	uniqueFilename := fmt.Sprintf("%d_%s", timestamp, filename)
	filePath := filepath.Join(s.basePath, uniqueFilename)

	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return uniqueFilename, nil
}

func (s *LocalFileStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, path)
}

func (s *LocalFileStorage) Delete(ctx context.Context, path string) error {
	filePath := filepath.Join(s.basePath, path)
	return os.Remove(filePath)
}
