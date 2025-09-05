package bookmarks

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Bookmarks struct {
	Named   map[string]string
	Unnamed map[string]string
}

var ErrBookmarkNotFound = errors.New("bookmark not found")

func NewBookmarks() *Bookmarks {
	return &Bookmarks{
		Named:   map[string]string{},
		Unnamed: map[string]string{},
	}
}

func (bm Bookmarks) AddNamed(name string, path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("cannot get absolut path: %q", err)
	}

	bm.Named[name] = absolutePath
	return nil
}

func (bm Bookmarks) GetNamed(name string) (string, error) {
	found, ok := bm.Named[name]

	if !ok {
		return "", fmt.Errorf("%w: %q", ErrBookmarkNotFound, name)
	}
	return found, nil
}

func (bm Bookmarks) PrettyList() string {
	maxNameLength := 0
	maxPathLength := 0
	for name, path := range bm.Named {
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
		if len(path) > maxPathLength {
			maxPathLength = len(path)
		}
	}

	output := ""
	for name, path := range bm.Named {
		output += fmt.Sprintf("| %-*s | %-*s |\n", maxNameLength, name, maxPathLength, path)
	}
	output = strings.TrimSuffix(output, "\n")

	return output
}

func (bm Bookmarks) RemoveNamed(name string) error {
	_, err := bm.GetNamed(name)
	if errors.Is(err, ErrBookmarkNotFound) {
		return err
	}
	delete(bm.Named, name)
	return nil
}
