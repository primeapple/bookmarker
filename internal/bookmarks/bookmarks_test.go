package bookmarks

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetNamed(t *testing.T) {
	bm := createBookmarks(map[string][]string{"home": {"/home/user"}})

	t.Run("find existing bookmark", func(t *testing.T) {
		want := []string{"/home/user"}

		got, err := bm.GetNamed("home")

		assertError(t, err, nil)
		assertPaths(t, got, want)
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

		err := bm.AddNamed(name, path)

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{path})
	})

	t.Run("new bookmark with multiple paths", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.AddNamed(name, "/b/second", "/a/first")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first", "/b/second"})
	})

	t.Run("overwrite existing bookmark", func(t *testing.T) {
		newPath := "/home/otherUser"
		bm := createBookmarks(map[string][]string{name: {path}})

		err := bm.AddNamed(name, newPath)

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{newPath})
	})

	t.Run("deduplicate paths", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.AddNamed(name, path, path)

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{path})
	})

	t.Run("make paths absolute", func(t *testing.T) {
		bm := *NewBookmarks()
		absolute, absErr := filepath.Abs("relative/dir")
		assertNil(t, absErr)

		err := bm.AddNamed(name, "relative/dir")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{absolute})
	})

	t.Run("give error without paths", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.AddNamed(name)

		assertNotNil(t, err)
	})
}

func TestAttachNamed(t *testing.T) {
	name := "ui"

	t.Run("attach to existing bookmark", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.AttachNamed(name, "/b/second")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first", "/b/second"})
	})

	t.Run("attach multiple paths at once", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.AttachNamed(name, "/c/third", "/b/second")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first", "/b/second", "/c/third"})
	})

	t.Run("create bookmark when it does not exist", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.AttachNamed(name, "/a/first")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first"})
	})

	t.Run("attaching an already attached path is a no op", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.AttachNamed(name, "/a/first")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first"})
	})

	t.Run("give error without paths", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.AttachNamed(name)

		assertNotNil(t, err)
	})
}

func TestDetachNamed(t *testing.T) {
	name := "ui"

	t.Run("detach one of multiple paths", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first", "/b/second"}})

		err := bm.DetachNamed(name, "/b/second")

		assertNil(t, err)
		assertBookmark(t, bm, name, []string{"/a/first"})
	})

	t.Run("remove bookmark when last path is detached", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.DetachNamed(name, "/a/first")

		assertNil(t, err)

		_, err = bm.GetNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("give error on non existing bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.DetachNamed(name, "/a/first")

		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("give error and change nothing on non attached path", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first", "/b/second"}})

		err := bm.DetachNamed(name, "/a/first", "/unknown")

		assertError(t, err, ErrPathNotAttached)
		assertBookmark(t, bm, name, []string{"/a/first", "/b/second"})
	})

	t.Run("give error without paths", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {"/a/first"}})

		err := bm.DetachNamed(name)

		assertNotNil(t, err)
	})
}

func TestPrettyList(t *testing.T) {
	t.Run("should pad with spaces and sort correctly", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{"name": {"path1"}, "verylongname": {"path2"}, "a_name": {"path3"}})
		want :=
			`| a_name       | path3 |
| name         | path1 |
| verylongname | path2 |`

		got := bm.PrettyList()
		assertString(t, got, want)
	})

	t.Run("should list every path of a bookmark on its own row", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{"ui": {"/a/first", "/b/second"}, "api": {"/c/third"}})
		want :=
			`| api | /c/third  |
| ui  | /a/first  |
|     | /b/second |`

		got := bm.PrettyList()
		assertString(t, got, want)
	})

	t.Run("should be empty without bookmarks", func(t *testing.T) {
		bm := *NewBookmarks()

		got := bm.PrettyList()
		assertString(t, got, "")
	})
}

func TestPorcelainList(t *testing.T) {
	t.Run("should print one tab separated row per path", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{"ui": {"/a/first", "/b/second"}, "api": {"/c/third"}})
		want := "api\t/c/third\nui\t/a/first\nui\t/b/second"

		got := bm.PorcelainList()
		assertString(t, got, want)
	})

	t.Run("should be empty without bookmarks", func(t *testing.T) {
		bm := *NewBookmarks()

		got := bm.PorcelainList()
		assertString(t, got, "")
	})
}

func TestRemoveNamed(t *testing.T) {
	name := "home"
	path := "/home/user"

	t.Run("remove existing bookmark", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {path}})

		err := bm.RemoveNamed(name)

		assertNil(t, err)

		_, err = bm.GetNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("remove multiple bookmarks", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {path}, "other": {"/other"}})

		err := bm.RemoveNamed(name, "other")

		assertNil(t, err)

		_, err = bm.GetNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
		_, err = bm.GetNamed("other")
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("remove non existing bookmark", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.RemoveNamed(name)
		assertError(t, err, ErrBookmarkNotFound)
	})

	t.Run("change nothing when one of the bookmarks does not exist", func(t *testing.T) {
		bm := createBookmarks(map[string][]string{name: {path}})

		err := bm.RemoveNamed(name, "unknown")

		assertError(t, err, ErrBookmarkNotFound)
		assertBookmark(t, bm, name, []string{path})
	})

	t.Run("give error without names", func(t *testing.T) {
		bm := *NewBookmarks()

		err := bm.RemoveNamed()

		assertNotNil(t, err)
	})
}

func createBookmarks(named map[string][]string) Bookmarks {
	return Bookmarks{Named: named}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if !errors.Is(got, want) {
		t.Errorf("got %v wanted %v", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q wanted %q", got, want)
	}
}

func assertPaths(t testing.TB, got, want []string) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v wanted %v", got, want)
	}
}

func assertBookmark(t testing.TB, bm Bookmarks, name string, paths []string) {
	t.Helper()

	got, err := bm.GetNamed(name)

	if err != nil {
		t.Fatalf("Should find added bookmark %q", name)
	}

	assertPaths(t, got, paths)
}

func assertNil(t testing.TB, got any) {
	t.Helper()

	if got != nil {
		t.Errorf("got %v wanted nil", got)
	}
}

func assertNotNil(t testing.TB, got error) {
	t.Helper()

	if got == nil {
		t.Errorf("got nil wanted an error")
	}
}
