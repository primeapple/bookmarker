package storage

import (
	"encoding/json"
	"fmt"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

type Migrator func(map[string]any) (map[string]any, error)

var migrations = map[int]Migrator{
	0: MigrateV0toV1,
}

func ParseBookmarksFile(jsonData []byte) (*bookmarks.Bookmarks, error) {
	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse json data into map[string]any: %w", err)
	}

	version, ok := data["_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid `_version` property in json data")
	}

	for int(version) < LATEST_VERSION {
		currentVersion := int(version)
		migrate, exists := migrations[currentVersion]
		if !exists {
			return nil, fmt.Errorf("no migration found for version %d", currentVersion)
		}

		data, err := migrate(data)
		if err != nil {
			return nil, fmt.Errorf("migration %d failed: %w", currentVersion, err)
		}

		version, ok = data["_version"].(float64)
		if !ok {
			return nil, fmt.Errorf("missing or invalid `_version` property in json data after migration %d", currentVersion)
		}
	}

	var latest LatestSchema
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not convert migrated data to json: %w", err)
	}
	if err := json.Unmarshal(b, &latest); err != nil {
		return nil, fmt.Errorf("could not fit migrated data into latest scchema: %w", err)
	}
	return latest.ToBookmarks(), nil
}
