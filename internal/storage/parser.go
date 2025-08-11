package storage

import (
	"encoding/json"
	"fmt"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

type Parser interface {
	Parse(data []byte) (*bookmarks.Bookmarks, error)
	version() int
}

type VersionedData struct {
	Version int `json:"_version"`
}

func ParseBookmarksFile(data []byte) (*bookmarks.Bookmarks, error) {
	version, err := detectVersion(data)
	if err != nil {
		return nil, err
	}

	parser := getParser(version)
	if parser == nil {
		return nil, fmt.Errorf("can't find a parser for version: %d", version)
	}

	bookmarks, err := parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	return bookmarks, nil
}

func detectVersion(data []byte) (int, error) {
	var versionData VersionedData
	if err := json.Unmarshal(data, &versionData); err != nil {
		return 0, fmt.Errorf("failed to parse version: %w", err)
	}

	return versionData.Version, nil
}

func getParser(version int) Parser {
	switch version {
	case 1:
		return V1Parser{}
	default:
		return nil
	}
}
