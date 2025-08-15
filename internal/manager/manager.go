package manager

import (
	"fmt"
	"os"

	"github.com/primeapple/bookmarker/internal/shell"
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
	case "add":
		return m.handleAdd(args[1:])
	case "get":
		return m.handleGet(args[1:])
	case "help":
		m.handlePrintHelp()
		return nil
	case "init":
		return m.handleInit(args[1:])
	case "list":
		return m.handleList()
	case "remove":
		return m.handleRemove(args[1:])
	default:
		m.handlePrintHelp()
		return nil
	}
}

func (m *Manager) handleAdd(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("only one argument is allowed for add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	bm.AddNamed(args[0], workingDirectory)
	err = m.store.Save(bm)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) handleGet(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("only one argument is allowed for add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	path, err := bm.GetNamed(names[0])
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %q for bookmark %q doesn't exist on disk", path, names[0])
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

	fmt.Println(bm.PrettyList())
	return nil
}

func (m *Manager) handleRemove(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("only one argument is allowed for add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	err = bm.RemoveNamed(names[0])
	if err != nil {
		return err
	}

	err = m.store.Save(bm)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) handleInit(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("only one argument is allowed for init")
	}

	switch names[0] {
	case "fish":
		fmt.Println(shell.InitFish())
	default:
		return fmt.Errorf("unsupported shell %q", names[0])
	}

	return nil
}

func (m *Manager) handlePrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bookmarker add")
	fmt.Println("  bookmarker get <name>")
	fmt.Println("  bookmarker init {fish}")
	fmt.Println("  bookmarker help")
	fmt.Println("  bookmarker list")
	fmt.Println("  bookmarker remove <name>")
}
