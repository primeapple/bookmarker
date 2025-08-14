package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

const BookmarkerFilename = "bookmarker.json"
const PermissionUserReadWrite = 0600
const PersmissionUserAllRestRead = 0711

type Storage interface {
	Load() (*bookmarks.Bookmarks, error)
	Save(*bookmarks.Bookmarks) error
}

type JSONStorage struct{}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{}
}

func (store *JSONStorage) Load() (*bookmarks.Bookmarks, error) {
	path, err := getStorageFilePath()
	if err != nil {
		return nil, fmt.Errorf("getting storage file dir: %w", err)
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return bookmarks.NewBookmarks(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("can't read bookmarks file: %w", err)
	}

	bookmarks, err := ParseBookmarksData(data)
	if err != nil {
		return nil, fmt.Errorf("can't parse bookmarks file: %w", err)
	}

	return bookmarks, nil
}

func (store *JSONStorage) Save(bm *bookmarks.Bookmarks) error {
	if bm == nil {
		panic("Passed `nil` as bookmarks")
	}

	schema := FromBookmarks(bm)

	data, err := json.MarshalIndent(schema, "", "\t")
	if err != nil {
		return fmt.Errorf("could not convert bookmarks schema %v to json %w", schema, err)
	}

	if err := store.createDirIfNotExists(); err != nil {
		return fmt.Errorf("creating storage directory: %w", err)
	}

	path, err := getStorageFilePath()
	if err := os.WriteFile(path, data, PermissionUserReadWrite); err != nil {
		return fmt.Errorf("could not write bookmarker file to %q , %w", path, err)
	}

	return nil
}

func (store *JSONStorage) createDirIfNotExists() error {
	path, err := getStorageFileDir()
	if err != nil {
		return fmt.Errorf("getting storage file dir: %w", err)
	}

	err = os.MkdirAll(path, PersmissionUserAllRestRead)
	if err != nil {
		return fmt.Errorf("creating dir %q: %w", path, err)
	}
	return nil
}

func getStorageFileDir() (string, error) {
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, "bookmarker"), nil
	}

	baseDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("homedir not found: %w", err)
	}

	return filepath.Join(baseDir, ".local", "share", "bookmarker"), nil
}

func getStorageFilePath() (string, error) {
	baseDir, err := getStorageFileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, BookmarkerFilename), nil
}
