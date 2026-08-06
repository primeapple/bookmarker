package storage

import (
	"testing"
)

func TestMigrateV1toV2(t *testing.T) {
	t.Run("migrate valid v1Data", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named": map[string]any{
					"name":      "path",
					"otherName": "otherPath",
				},
				"unnamed": map[string]any{},
			},
		}
		want := map[string]any{
			"_version": 2,
			"bookmarks": map[string]any{
				"named": map[string][]string{
					"name":      {"path"},
					"otherName": {"otherPath"},
				},
			},
		}

		got, err := MigrateV1toV2(v1Data)

		assertNil(t, err)
		assertMap(t, got, want)
	})

	t.Run("migrate empty v1Data", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named":   map[string]any{},
				"unnamed": map[string]any{},
			},
		}
		want := map[string]any{
			"_version": 2,
			"bookmarks": map[string]any{
				"named": map[string][]string{},
			},
		}

		got, err := MigrateV1toV2(v1Data)

		assertNil(t, err)
		assertMap(t, got, want)
	})

	t.Run("abort on wrong version", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 0,
			"bookmarks": map[string]any{
				"named":   map[string]any{},
				"unnamed": map[string]any{},
			},
		}

		_, err := MigrateV1toV2(v1Data)

		assertNotNil(t, err)
	})

	t.Run("abort on missing version", func(t *testing.T) {
		v1Data := map[string]any{
			"bookmarks": map[string]any{
				"named":   map[string]any{},
				"unnamed": map[string]any{},
			},
		}

		_, err := MigrateV1toV2(v1Data)

		assertNotNil(t, err)
	})

	t.Run("abort on missing bookmarks", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 1,
		}

		_, err := MigrateV1toV2(v1Data)

		assertNotNil(t, err)
	})

	t.Run("abort on non string path", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named":   map[string]any{"name": 42},
				"unnamed": map[string]any{},
			},
		}

		_, err := MigrateV1toV2(v1Data)

		assertNotNil(t, err)
	})

	t.Run("abort on non empty unnamed data", func(t *testing.T) {
		v1Data := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named":   map[string]any{},
				"unnamed": map[string]any{"name": "path"},
			},
		}

		_, err := MigrateV1toV2(v1Data)

		assertNotNil(t, err)
	})
}
