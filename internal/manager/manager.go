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
	attachCmd := flag.NewFlagSet("attach", flag.ExitOnError)
	detachCmd := flag.NewFlagSet("detach", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)

	listPorcelain := listCmd.Bool("porcelain", false, "print machine readable `name path1 path2 ... pathN` lines instead of a table")

	addCmd.Usage = func() {
		fmt.Printf("Usage: %s add [name] path...\n", os.Args[0])
		fmt.Printf("Add a named or temporary item to the bookmarks, replacing all its paths\n")
	}
	attachCmd.Usage = func() {
		fmt.Printf("Usage: %s attach name path...\n", os.Args[0])
		fmt.Printf("Attach additional paths to a bookmark, creating it if needed\n")
	}
	detachCmd.Usage = func() {
		fmt.Printf("Usage: %s detach name path...\n", os.Args[0])
		fmt.Printf("Detach paths from a bookmark, removing it if no path is left\n")
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
	case "attach":
		err := attachCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleAttach(attachCmd.Args())
	case "detach":
		err := detachCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleDetach(detachCmd.Args())
	case "get":
		err := getCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleGet(getCmd.Args())
	case "go":
		err := goCmd.Parse(args[1:])
		if err != nil {
			return err
		}
		return m.handleGo(goCmd.Args())
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
		return m.handleList(*listPorcelain)
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
	if len(args) == 0 {
		return fmt.Errorf("at least one argument is required for add")
	}

	if len(args) == 1 {
		return fmt.Errorf("temporary bookmarks not yet implemented")
	}

	name := args[0]
	paths := args[1:]

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if err := bm.AddNamed(name, paths...); err != nil {
		return fmt.Errorf("bookmark cannot be added: %w", err)
	}

	return m.store.Save(bm)
}

func (m *Manager) handleAttach(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("a name and at least one path are required for attach")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if err := bm.AttachNamed(args[0], args[1:]...); err != nil {
		return fmt.Errorf("paths cannot be attached: %w", err)
	}

	return m.store.Save(bm)
}

func (m *Manager) handleDetach(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("a name and at least one path are required for detach")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if err := bm.DetachNamed(args[0], args[1:]...); err != nil {
		return fmt.Errorf("paths cannot be detached: %w", err)
	}

	return m.store.Save(bm)
}

func (m *Manager) handleGet(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("only one argument is allowed for get")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	paths, err := bm.GetNamed(names[0])
	if err != nil {
		return err
	}

	for _, path := range paths {
		fmt.Println(path)
	}
	return nil
}

func (m *Manager) handleGo(names []string) error {
	if len(names) != 1 {
		return fmt.Errorf("only one argument is allowed for go")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	paths, err := bm.GetNamed(names[0])
	if err != nil {
		return err
	}

	if len(paths) != 1 {
		return fmt.Errorf("bookmark %q holds %d paths, cannot go to more than one directory", names[0], len(paths))
	}

	path := paths[0]
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path %q for bookmark %q doesn't exist on disk", path, names[0])
		}
		return err
	}

	fmt.Println(path)
	return nil
}

func (m *Manager) handleList(porcelain bool) error {
	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if porcelain {
		printIfNotEmpty(bm.PorcelainList())
		return nil
	}

	printIfNotEmpty(bm.PrettyList())
	return nil
}

func printIfNotEmpty(output string) {
	if output == "" {
		return
	}
	fmt.Println(output)
}

func (m *Manager) handleRemove(names []string) error {
	if len(names) == 0 {
		return fmt.Errorf("at least one argument is required for remove")
	}

	bm, err := m.store.Load()
	if err != nil {
		return err
	}

	if err := bm.RemoveNamed(names...); err != nil {
		return err
	}

	return m.store.Save(bm)
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
	fmt.Printf("  %s add <name> <path>...\n", os.Args[0])
	fmt.Printf("  %s attach <name> <path>...\n", os.Args[0])
	fmt.Printf("  %s detach <name> <path>...\n", os.Args[0])
	fmt.Printf("  %s get <name>\n", os.Args[0])
	fmt.Printf("  %s go <name>\n", os.Args[0])
	fmt.Printf("  %s help\n", os.Args[0])
	fmt.Printf("  %s init {fish}\n", os.Args[0])
	fmt.Printf("  %s list [--porcelain]\n", os.Args[0])
	fmt.Printf("  %s remove <name>...\n", os.Args[0])
}
