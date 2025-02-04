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
		// return m.printUsage()
	}

	switch args[0] {
	case "--add":
		return m.handleAdd(args[1:])
	default:
		return m.handleGoto(args[0])
	}
}

func (m *Manager) handleAdd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Wrong number of arguments given")
	}

	bm := m.store.Load()

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic("Can't get current working directory")
	}

	bm.Add(args[0], workingDirectory)
	return nil
}

func (m *Manager) handleGoto(name string) error {
	bm := m.store.Load()
    path, err := bm.Get(name)
    if err != nil {
        return err
    }

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Path %q for bookmark %q doesn't exist", path, name)
		} else {
			panic(err)
		}
	}

    println(path)

	return nil
}
