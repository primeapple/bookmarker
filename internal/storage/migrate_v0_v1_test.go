package storage

import (
	"reflect"
	"testing"
)

func TestMigrateV0toV1(t *testing.T) {
	t.Run("migrate valid v0Data", func(t *testing.T) {
		v0Data := map[string]any{
			"_version":  0,
			"name":      "path",
			"otherName": "otherPath",
		}
		want := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named": map[string]any{
					"name":      "path",
					"otherName": "otherPath",
				},
				"unnamed": map[string]any{},
			},
		}

		got, err := MigrateV0toV1(v0Data)

		assertNil(t, err)
		assertMap(t, got, want)
	})

	t.Run("migrate empty v0Data", func(t *testing.T) {
		v0Data := map[string]any{
			"_version": 0,
		}
		want := map[string]any{
			"_version": 1,
			"bookmarks": map[string]any{
				"named":   map[string]any{},
				"unnamed": map[string]any{},
			},
		}

		got, err := MigrateV0toV1(v0Data)

		assertNil(t, err)
		assertMap(t, got, want)
	})

	t.Run("abort on wrong version", func(t *testing.T) {
		v0Data := map[string]any{
			"_version":  1,
			"name":      "path",
			"otherName": "otherPath",
		}

		_, err := MigrateV0toV1(v0Data)

		assertNotNil(t, err)
	})

	t.Run("abort on missing version", func(t *testing.T) {
		v0Data := map[string]any{
			"name":      "path",
			"otherName": "otherPath",
		}

		_, err := MigrateV0toV1(v0Data)

		assertNotNil(t, err)
	})
}

func assertMap(t testing.TB, got, want map[string]any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
