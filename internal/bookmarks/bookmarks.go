package bookmarks

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

type Bookmarks struct {
	Named map[string][]string
}

var ErrBookmarkNotFound = errors.New("bookmark not found")
var ErrPathNotAttached = errors.New("path not attached to bookmark")

func NewBookmarks() *Bookmarks {
	return &Bookmarks{
		Named: map[string][]string{},
	}
}

func (bm Bookmarks) AddNamed(name string, paths ...string) error {
	normalized, err := normalizePaths(paths)
	if err != nil {
		return err
	}

	bm.Named[name] = normalized
	return nil
}

func (bm Bookmarks) AttachNamed(name string, paths ...string) error {
	normalized, err := normalizePaths(paths)
	if err != nil {
		return err
	}

	bm.Named[name] = normalizePathSet(append(bm.Named[name], normalized...))
	return nil
}

func (bm Bookmarks) DetachNamed(name string, paths ...string) error {
	current, err := bm.GetNamed(name)
	if err != nil {
		return err
	}

	normalized, err := normalizePaths(paths)
	if err != nil {
		return err
	}

	for _, path := range normalized {
		if !slices.Contains(current, path) {
			return fmt.Errorf("%w: %q is not attached to %q", ErrPathNotAttached, path, name)
		}
	}

	remaining := slices.DeleteFunc(slices.Clone(current), func(path string) bool {
		return slices.Contains(normalized, path)
	})

	if len(remaining) == 0 {
		delete(bm.Named, name)
		return nil
	}

	bm.Named[name] = remaining
	return nil
}

func (bm Bookmarks) GetNamed(name string) ([]string, error) {
	found, ok := bm.Named[name]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrBookmarkNotFound, name)
	}
	return found, nil
}

func (bm Bookmarks) RemoveNamed(names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("at least one name is required")
	}

	for _, name := range names {
		if _, err := bm.GetNamed(name); err != nil {
			return err
		}
	}

	for _, name := range names {
		delete(bm.Named, name)
	}
	return nil
}

func (bm Bookmarks) PrettyList() string {
	maxNameLength := 0
	maxPathLength := 0
	for name, paths := range bm.Named {
		if len(name) > maxNameLength {
			maxNameLength = len(name)
		}
		for _, path := range paths {
			if len(path) > maxPathLength {
				maxPathLength = len(path)
			}
		}
	}

	output := ""
	for _, name := range bm.sortedNames() {
		for index, path := range bm.Named[name] {
			displayedName := name
			if index > 0 {
				displayedName = ""
			}
			output += fmt.Sprintf("| %-*s | %-*s |\n", maxNameLength, displayedName, maxPathLength, path)
		}
	}

	return strings.TrimSuffix(output, "\n")
}

func (bm Bookmarks) PorcelainList() string {
	output := ""
	for _, name := range bm.sortedNames() {
		for _, path := range bm.Named[name] {
			output += fmt.Sprintf("%s\t%s\n", name, path)
		}
	}

	return strings.TrimSuffix(output, "\n")
}

func (bm Bookmarks) sortedNames() []string {
	names := make([]string, 0, len(bm.Named))
	for name := range bm.Named {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func normalizePaths(paths []string) ([]string, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("at least one path is required")
	}

	absolutePaths := make([]string, 0, len(paths))
	for _, path := range paths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("cannot get absolute path for %q: %w", path, err)
		}
		absolutePaths = append(absolutePaths, absolutePath)
	}

	return normalizePathSet(absolutePaths), nil
}

func normalizePathSet(paths []string) []string {
	slices.Sort(paths)
	return slices.Compact(paths)
}

// PathType represents the type of a filesystem path
type PathType string

const (
	FileType PathType = "file"
	DirType  PathType = "dir"
)

func typeOf(path string) (PathType, error) {

}
