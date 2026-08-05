package storage

import "fmt"

const V1_VERSION = 1

func MigrateV1toV2(data map[string]any) (map[string]any, error) {
	version, ok := data["_version"].(int)
	if !ok {
		return nil, fmt.Errorf("missing or invalid `_version` property in json data: %v", data)
	}

	if version != V1_VERSION {
		return nil, fmt.Errorf("wrong versioned file for V1toV2 migration, got version %d", version)
	}

	oldBookmarks, ok := data["bookmarks"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid `bookmarks` property in json data: %v", data)
	}

	unnamed, ok := oldBookmarks["unnamed"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'unnamed' property in bookmarks: %v", oldBookmarks)
	}
	if len(unnamed) != 0 {
		return nil, fmt.Errorf("there are entries in the 'unnamed' property in bookmarks, but there shouldn't: %v", oldBookmarks)
	}

	oldNamed, ok := oldBookmarks["named"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'named' property in bookmarks: %v", oldBookmarks)
	}

	named := map[string][]string{}
	for name, path := range oldNamed {
		stringPath, ok := path.(string)
		if !ok {
			return nil, fmt.Errorf("invalid non string path for bookmark %q in bookmarks: %v", name, oldBookmarks)
		}

		named[name] = []string{stringPath}
	}

	bookmarks := map[string]any{
		"named": named,
	}

	result := map[string]any{
		"_version":  V1_VERSION + 1,
		"bookmarks": bookmarks,
	}

	return result, nil
}
