package bookmarks

import (
	"errors"
	"testing"
)

func TestGetNamed(t *testing.T) {
	bm := createBookmarks(map[string]string{"home": "/home/user"})

	t.Run("find existing bookmark", func(t *testing.T) {
		want := "/home/user"

		got, err := bm.GetNamed("home")

		assertError(t, err, nil)
		assertString(t, got, want)
	})

	t.Run("give error on non existing bookmark", func(t *testing.T) {
		_, err := bm.GetNamed("unknown")

		assertError(t, err, ErrBookmarkNotFound)
	})
}

func TestAddNamed(t *testing.T) {
	name := "home"
	path := "/home/user"

	t.Run("new bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

		bm.AddNamed(name, path)

		assertBookmark(t, bm, name, path)
	})

	t.Run("overwrite existing bookmark", func(t *testing.T) {
		newPath := "/home/otherUser"
		bm := createBookmarks(map[string]string{name: path})

		bm.AddNamed(name, newPath)

		assertBookmark(t, bm, name, newPath)
	})
}

func TestPrettyList(t *testing.T) {
	bm := createBookmarks(map[string]string{"name": "path1", "verylongname": "path2"})

	t.Run("should pad with spaces correctly", func(t *testing.T) {
		want :=
			`| name         | path1 |
| verylongname | path2 |`

		got := bm.PrettyList()
		assertString(t, got, want)
	})
}

func TestRemoveNamed(t *testing.T) {
	name := "home"
	path := "/home/user"

	t.Run("remove existing bookmark", func(t *testing.T) {
		bm := createBookmarks(map[string]string{name: path})

		bm.RemoveNamed(name)

		_, err := bm.GetNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("remove non existing bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.RemoveNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
	})
}

func createBookmarks(named map[string]string) Bookmarks {
	return Bookmarks{Named: named, Unnamed: map[string]string{}}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("got %q wanted %q", got.Error(), want.Error())
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q wanted %q", got, want)
	}
}

func assertBookmark(t testing.TB, bm Bookmarks, name, path string) {
	got, err := bm.GetNamed(name)

	if err != nil {
		t.Fatalf("Should find added bookmark %q", name)
	}

	assertString(t, got, path)
}
