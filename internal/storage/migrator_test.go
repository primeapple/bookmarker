package storage

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/primeapple/bookmarker/internal/bookmarks"
)

func TestParseBookmarksData(t *testing.T) {
	t.Run("migrate valid v0Data", func(t *testing.T) {
		v0Data := []byte(`{
			"_version":  0,
			"name":      "path",
			"otherName": "otherPath"
		}`)
		want := bookmarks.Bookmarks{
			Named: map[string][]string{
				"name":      {"path"},
				"otherName": {"otherPath"},
			},
		}

		got, err := ParseBookmarksData(v0Data)

		assertNil(t, err)
		assertBookmarks(t, got, &want)
	})

	t.Run("migrate valid v1Data", func(t *testing.T) {
		v1Data := []byte(`{
			"_version":  1,
			"bookmarks": {
				"named": {
					"name":      "path",
					"otherName": "otherPath"
				},
				"unnamed": {}
			}
		}`)
		want := bookmarks.Bookmarks{
			Named: map[string][]string{
				"name":      {"path"},
				"otherName": {"otherPath"},
			},
		}

		got, err := ParseBookmarksData(v1Data)

		assertNil(t, err)
		assertBookmarks(t, got, &want)
	})

	t.Run("load valid latest data", func(t *testing.T) {
		latestData := fmt.Appendf(nil, `{
			"_version":  %d,
			"bookmarks": {
				"named": {
					"name":      ["path"],
					"otherName": ["otherPath", "secondPath"]
				},
				"unnamed": {}
			}
		}`, LATEST_VERSION)
		want := bookmarks.Bookmarks{
			Named: map[string][]string{
				"name":      {"path"},
				"otherName": {"otherPath", "secondPath"},
			},
		}

		got, err := ParseBookmarksData(latestData)

		assertNil(t, err)
		assertBookmarks(t, got, &want)
	})

	t.Run("abort on non existing version", func(t *testing.T) {
		data := fmt.Appendf(nil, `{
			"_version":  %d,
			"name":      "path",
			"otherName": "otherPath"
		}`, LATEST_VERSION+1)

		_, err := ParseBookmarksData(data)

		assertNotNil(t, err)
	})

	t.Run("abort on missing version", func(t *testing.T) {
		data := []byte(`{
			"name":      "path",
			"otherName": "otherPath"
		}`)

		_, err := ParseBookmarksData(data)

		assertNotNil(t, err)
	})
}

func assertNil(t testing.TB, got any) {
	if got != nil {
		t.Errorf("got %v wanted nil", got)
	}
}

func assertNotNil(t testing.TB, got error) {
	t.Helper()

	if got == nil {
		t.Errorf("got nil wanted different")
	}
}

func assertBookmarks(t testing.TB, got, want *bookmarks.Bookmarks) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
