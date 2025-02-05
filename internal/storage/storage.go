package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

const BOOKMARKER_FILENAME = "bookmarker.json"

type Storage interface {
	Load() *bookmarks.Bookmarks
	Save(*bookmarks.Bookmarks)
}

type JSONStorage struct{}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{}
}

func (store *JSONStorage) Load() *bookmarks.Bookmarks {
	path := getStorageFilePath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		newBookmarks := bookmarks.NewBookmarks()
		return newBookmarks
	}
	if err != nil {
		panic(fmt.Sprintf("can't read bookmarks file: %w", err))
	}

	var result bookmarks.Bookmarks
	if err := json.Unmarshal(data, &result); err != nil {
		panic(fmt.Sprintf("File %v was not a valid json: %w", data, err))
	}

	return &result
}

func (store *JSONStorage) Save(bm *bookmarks.Bookmarks) {
	data, err := json.MarshalIndent(bm, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Could not convert bookmarks %v to json %w", bm, err))
	}

    store.createDirIfNotExists()

    path := getStorageFilePath()
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		panic(fmt.Sprintf("Could not write bookmarker file to %q , %w", path, err))
	}
}

func (store *JSONStorage) createDirIfNotExists() {
	path := getStorageFileDir()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0711)
	}
}

func getStorageFileDir() string {
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, "bookmarker")
	}

	baseDir, err := os.UserHomeDir()
	if err != nil {
		baseDir = os.Getenv("HOME")
	}

	return filepath.Join(baseDir, ".local", "share", "bookmarker")
}

func getStorageFilePath() string {
	baseDir := getStorageFileDir()
	return filepath.Join(baseDir, BOOKMARKER_FILENAME)
}
