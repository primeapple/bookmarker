package manager

import (
	"fmt"
	"os"

	"github.com/primeapple/bookmarker/internal/storage"
)

type Manager struct {
	store storage.Storage
}

func NewManager() *Manager {
	return &Manager{
		store: storage.NewJSONStorage(),
	}
}

func (m *Manager) Run(args []string) error {
	if len(args) == 0 {
		m.handlePrintHelp()
        return nil
	}

	switch args[0] {
	case "-a", "--add":
		return m.handleAdd(args[1:])
	case "-g", "--get":
		return m.handleGet(args[1])
    case "-l", "--list":
		return m.handleList()
	case "-h", "--help":
		m.handlePrintHelp()
        return nil
	default:
		m.handlePrintHelp()
        return nil
	}
}

func (m *Manager) handleAdd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Only one argument is allowed for --add")
	}

	bm, err := m.store.Load()
	if err != nil {
        return err
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	bm.Add(args[0], workingDirectory)
    err = m.store.Save(bm)
    if err != nil {
        return err
    }

    return nil
}

func (m *Manager) handleGet(name string) error {
	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	path, err := bm.Get(name)
	if err != nil {
        return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
            return fmt.Errorf("Path %q for bookmark %q doesn't exist on disk", path, name)
		} else {
			return err
		}
	}

	fmt.Println(path)
    return nil
}

func (m *Manager) handleList() error {
	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	bookmarks := bm.ListAll()
    for i, bookmark := range bookmarks {
        fmt.Println(i, bookmark.Name, bookmark.Path)
    }
    return nil
}

func (m *Manager) handlePrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bookmarker [-a | --add] <name>")
	fmt.Println("  bookmarker [-g | --get] <name>")
	fmt.Println("  bookmarker [-h | --help]")
	fmt.Println("  bookmarker [-l | --list]")
}
