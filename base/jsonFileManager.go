package base

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type FileManager[T any] struct {
	basePath string
	mu       sync.RWMutex
}

func NewJsonFileManager[T any](basePath string) (*FileManager[T], error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, ErrFailedToCreateDirectory
	}

	return &FileManager[T]{
		basePath: basePath,
	}, nil
}

func (fm *FileManager[T]) Write(filename string, data T) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	fullPath := filepath.Join(fm.basePath, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ErrFailedJSONParse
	}

	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return ErrFailedToWriteFile
	}

	return nil
}

func (fm *FileManager[T]) Read(filename string) (T, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	fullPath := filepath.Join(fm.basePath, filename)

	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		return *new(T), ErrFileNotExists
	}

	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return *new(T), ErrFailedToReadFile
	}

	var result T
	if err := json.Unmarshal(fileData, &result); err != nil {
		return *new(T), ErrFailedJSONParse
	}

	return result, nil
}

func (fm *FileManager[T]) Exists(filename string) bool {
	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	fullPath := filepath.Join(fm.basePath, filename)
	_, err := os.Stat(fullPath)
	return !errors.Is(err, os.ErrNotExist)
}

func (fm *FileManager[T]) Delete(filename string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	fullPath := filepath.Join(fm.basePath, filename)
	return os.Remove(fullPath)
}
