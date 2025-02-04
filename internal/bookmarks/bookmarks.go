package bookmarks

import (
	"fmt"
)

type Bookmarks map[string]string

func NewBookmarks() *Bookmarks {
	return &Bookmarks{}
}

func (bm Bookmarks) Add(name string, path string) {
	bm[name] = path
}

func (bm Bookmarks) Get(name string) (string, error) {
	found, ok := bm[name]

	if !ok {
		return "", fmt.Errorf("Could not find %q in bookmarks", name)
	}
	return found, nil
}
