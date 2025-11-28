package manager

import (
	"flag"
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
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)

	addCmd.Usage = func() {
		fmt.Printf("Usage: %s add [name] path\n", os.Args[0])
		fmt.Printf("Add a named or temporary item to the bookmarks\n")
	}

	if len(args) == 0 {
		m.handlePrintHelp()
		return nil
	}

	switch args[0] {
	case "add":
		err := addCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleAdd(addCmd.Args())
	case "get":
		err := getCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleGet(getCmd.Args())
	case "help":
		m.handlePrintHelp()
		return nil
	case "init":
		err := initCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleInit(initCmd.Args())
	case "list":
		err := listCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleList()
	case "remove":
		err := removeCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleRemove(removeCmd.Args())
	default:
		m.handlePrintHelp()
		return nil
	}
}

func (m *Manager) handleAdd(args []string) error {
	var name string
	var path string
	switch len(args) {
	case 1:
		path = args[0]
	case 2:
		name = args[0]
		path = args[1]
	default:
		return fmt.Errorf("either one or two arguments are allowed for add")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if name == "" {
		// add temp bookmark
	} else {
		err := bm.AddNamed(name, path)
		if err != nil {
			return fmt.Errorf("bookmark cannot be added: %q", err)
		}
	}

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
	fmt.Printf("  %s add [<name>] <path>\n", os.Args[0])
	fmt.Printf("  %s get <name>\n", os.Args[0])
	fmt.Printf("  %s go <name>\n", os.Args[0])
	fmt.Printf("  %s init {fish}\n", os.Args[0])
	fmt.Printf("  %s help\n", os.Args[0])
	fmt.Printf("  %s list\n", os.Args[0])
	fmt.Printf("  %s remove <name>\n", os.Args[0])
}
