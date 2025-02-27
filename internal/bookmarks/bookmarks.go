package bookmarks

import (
	"errors"
	"fmt"
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

func (bm Bookmarks) Remove(name string) error {
    _, err := bm.Get(name)
    if errors.Is(err, ErrBookmarkNotFound) {
        return err 
    }
    delete(bm, name)
    return nil
}
