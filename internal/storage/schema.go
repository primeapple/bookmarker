package storage

import "github.com/primeapple/bookmarker/internal/bookmarks"

const LATEST_VERSION = 1

type LatestSchema struct {
	Version   int `json:"_version"`
	Bookmarks struct {
		Named   map[string]string `json:"named"`
		Unnamed map[string]string `json:"unnamed"`
	} `json:"bookmarks"`
}

func FromBookmarks(b *bookmarks.Bookmarks) *LatestSchema {
	schema := &LatestSchema{
		Version: LATEST_VERSION,
	}
	schema.Bookmarks.Named = b.Named
	schema.Bookmarks.Unnamed = b.Unnamed
	return schema
}

func (schema *LatestSchema) ToBookmarks() *bookmarks.Bookmarks {
	return &bookmarks.Bookmarks{
		Named:   schema.Bookmarks.Named,
		Unnamed: schema.Bookmarks.Unnamed,
	}
}
