package storage

import "fmt"

const BEFORE_VERSION = 0

func MigrateV0toV1(data map[string]any) (map[string]any, error) {
	version, ok := data["_version"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid `_version` property in json data")
	}
	if int(version) != BEFORE_VERSION {
		return nil, fmt.Errorf("wrong versioned file for V0toV1 migration, got version %d", int(version))
	}

	named := map[string]any{}

	for name, path := range data {
		if name == "_version" {
			continue
		}

		named[name] = path
	}

	result := map[string]any{
		"_version": BEFORE_VERSION + 1,
		"named":    named,
		"unnamed":  map[string]any{},
	}

	return result, nil
}
