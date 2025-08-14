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

func ParseBookmarksData(jsonData []byte) (*bookmarks.Bookmarks, error) {
	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse json data into map[string]any: %w", err)
	}

	version, ok := data["_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid `_version` property in json data %v", data)
	}
	intVersion := int(version)
	data["_version"] = intVersion

	_, exists := migrations[intVersion]
	if !exists && intVersion != LATEST_VERSION {
		return nil, fmt.Errorf("no migration found for version %d", intVersion)
	}

	for intVersion < LATEST_VERSION {
		migrate, exists := migrations[intVersion]
		if !exists {
			return nil, fmt.Errorf("no migration found for version %d", intVersion)
		}

		var err error
		data, err = migrate(data)
		if err != nil {
			return nil, fmt.Errorf("migration %d failed: %w", int(version), err)
		}

		intVersion, ok = data["_version"].(int)
		if !ok {
			return nil, fmt.Errorf("missing or invalid `_version` property in json data after migration %d", intVersion)
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

func getVersion(data map[string]any) (float64, bool) {
	switch v := data["_version"].(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	default:
		return -1, false
	}
}
