package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

type Storage interface {
	Load() *bookmarks.Bookmarks
	Save(*bookmarks.Bookmarks)
}

type JSONStorage struct{}

func NewJSONStorage() *JSONStorage {
	return &JSONStorage{}
}

func (JSONStorage) Load() *bookmarks.Bookmarks {
	path := getStorageFilePath("bookmarker.json")
	file, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		newBookmarks := bookmarks.NewBookmarks()
		return newBookmarks
	}
	if err != nil {
		panic(fmt.Sprintf("can't read bookmarks file: %w", err))
	}

	var result bookmarks.Bookmarks
	if err := json.Unmarshal(file, &result); err != nil {
		panic(fmt.Sprintf("File `bookmarker.json` was not a valid json", err))
	}

	return &result
}

func (JSONStorage) Save(bm *bookmarks.Bookmarks) {
	data, err := json.MarshalIndent(bm, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("Could not convert bookmarks %v to json %w", bm, err))
	}

	path := getStorageFilePath("bookmarker.json")
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		panic(fmt.Sprintf("Could not write bookmarker file to %q , %w", path, err))
	}
}

func getStorageFilePath(name string) string {
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return filepath.Join(dataHome, name)
	}

	baseDir, err := os.UserHomeDir()
	if err != nil {
		baseDir = os.Getenv("HOME")
	}

	return filepath.Join(baseDir, ".local", "share", name)
}
