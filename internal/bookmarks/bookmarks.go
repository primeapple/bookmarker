package bookmarks

import (
	"errors"
	"fmt"
	"strings"
)

type Bookmarks map[string]string

var ErrBookmarkNotFound = errors.New("bookmark not found")

func NewBookmarks() *Bookmarks {
	return &Bookmarks{}
}

func (bm Bookmarks) Add(name string, path string) {
	bm[name] = path
}

func (bm Bookmarks) Get(name string) (string, error) {
	found, ok := bm[name]

	if !ok {
		return "", fmt.Errorf("%w: %q", ErrBookmarkNotFound, name)
	}
	return found, nil
}

func (bm Bookmarks) PrettyList() string {
	maxNameLength := 0
	maxPathLength := 0
	for name, path := range bm {
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
		if len(path) > maxPathLength {
			maxPathLength = len(path)
		}
	}

	output := ""
	for name, path := range bm {
		output += fmt.Sprintf("| %-*s | %-*s |\n", maxNameLength, name, maxPathLength, path)
	}
	output = strings.TrimSuffix(output, "\n")

	return output
}

func (bm Bookmarks) Remove(name string) error {
	_, err := bm.Get(name)
	if errors.Is(err, ErrBookmarkNotFound) {
		return err
	}
	delete(bm, name)
	return nil
}
