package bookmarks

import (
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
	bm := Bookmarks{"home": "/home/user"}

	t.Run("find existing bookmark", func(t *testing.T) {
		want := "/home/user"

		got, err := bm.Get("home")

		assertError(t, err, nil)
		assertString(t, got, want)
	})

	t.Run("give error on non existing bookmark", func(t *testing.T) {
		_, err := bm.Get("unknown")

		assertError(t, err, ErrBookmarkNotFound)
	})
}

func TestAdd(t *testing.T) {
	name := "home"
	path := "/home/user"

	t.Run("new bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

        bm.Add(name, path)

		assertBookmark(t, bm, name, path)
	})

	t.Run("overwrite existing bookmark", func(t *testing.T) {
        newPath := "/home/otherUser"
		bm := Bookmarks{name: path}

		bm.Add(name, newPath)

		assertBookmark(t, bm, name, newPath)
	})
}

func TestRemove(t *testing.T) {
	name := "home"
	path := "/home/user"

	t.Run("remove existing bookmark", func(t *testing.T) {
		bm := Bookmarks{name: path}

		bm.Remove(name)

        _, err := bm.Get(name)
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("remove non existing bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

        err := bm.Remove(name)
		assertError(t, err, ErrBookmarkNotFound)
	})
}


func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("Got %q wanted %q", got.Error(), want.Error())
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("Got %q wanted %q", got, want)
	}
}

func assertBookmark(t testing.TB, bm Bookmarks, name, path string) {
	got, err := bm.Get(name)

	if err != nil {
		t.Fatalf("Should find added bookmark %q", name)
	}

	assertString(t, got, path)
}
