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
		return m.handleGet(args[1:])
	case "-l", "--list":
		return m.handleList()
	case "-h", "--help":
		m.handlePrintHelp()
		return nil
	case "-r", "--remove":
		return m.handleRemove(args[1:])
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

func (m *Manager) handleGet(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("Only one argument is allowed for --add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	path, err := bm.Get(names[0])
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Path %q for bookmark %q doesn't exist on disk", path, names[0])
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

	for name, path := range *bm {
		fmt.Println(name, path)
	}
	return nil
}

func (m *Manager) handleRemove(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("Only one argument is allowed for --add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

    err = bm.Remove(names[0])
    if err != nil {
        return err
    }

	err = m.store.Save(bm)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) handlePrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bookmarker [-a | --add] <name>")
	fmt.Println("  bookmarker [-g | --get] <name>")
	fmt.Println("  bookmarker [-h | --help]")
	fmt.Println("  bookmarker [-l | --list]")
	fmt.Println("  bookmarker [-r | --remove] <name>")
}
