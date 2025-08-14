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
			Named: map[string]string{
				"name":      "path",
				"otherName": "otherPath",
			},
			Unnamed: map[string]string{},
		}

		got, err := ParseBookmarksData(v0Data)

		assertNil(t, err)
		assertBookmarks(t, got, &want)
	})

	t.Run("abort on non existing version", func(t *testing.T) {
		data := []byte(`{
			"_version":  1,
			"name":      "path",
			"otherName": "otherPath"
		}`)

		_, err := ParseBookmarksData(data)

		print(fmt.Printf("%v", err))
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
