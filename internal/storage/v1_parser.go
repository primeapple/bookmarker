package storage

import (
	"encoding/json"
	"fmt"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

type V1Parser struct{}

func (p V1Parser) version() int {
	return 1
}

type V1Bookmarks struct {
	Version          int               `json:"_version"`
	NamedBookmarks   map[string]string `json:"namedBookmarks"`
	UnnamedBookmarks map[string]string `json:"unnamedBookmarks"`
}

func (p V1Parser) Parse(data []byte) (*bookmarks.Bookmarks, error) {
	var result V1Bookmarks
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("file %v is not a valid version 1 bookmarker file: %q", data, err.Error())
	}

	if result.Version != p.version() {
		return nil, fmt.Errorf("parser can only handle version %d files, but was given version %d", p.version(), result.Version)
	}


	return &bookmarks.Bookmarks{
		Named:   result.NamedBookmarks,
		Unnamed: result.UnnamedBookmarks,
	}, nil
}
