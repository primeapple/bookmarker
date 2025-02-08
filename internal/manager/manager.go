package manager

import (
	"fmt"
	"log"
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

func (m *Manager) Run(args []string) {
	if len(args) == 0 {
		m.handlePrintHelp()
	}

	switch args[0] {
	case "--add":
		m.handleAdd(args[1:])
	case "--get":
		m.handleGet(args[1])
	case "--help":
		m.handlePrintHelp()
	default:
		m.handlePrintHelp()
	}
}

func (m *Manager) handleAdd(args []string) {
	if len(args) != 1 {
		panic("Only one argument is allowed for --add")
	}

	bm, err := m.store.Load()
	if err != nil {
		log.Printf("error when loading the bookmarks: %v", err)
		panic("Unable to load the bookmarks. Please check permission and filespace")
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic("Can't get current working directory")
	}

	bm.Add(args[0], workingDirectory)
	m.store.Save(bm)

    os.Exit(0)
}

func (m *Manager) handleGet(name string) {
	bm, err := m.store.Load()
	if err != nil {
		log.Printf("error when loading the bookmarks: %v", err)
		panic("Unable to load the bookmarks. Please check permission and filespace")
	}

	path, err := bm.Get(name)
	if err != nil {
		panic(fmt.Sprintf("No bookmark found with name %q", name))
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
            panic(fmt.Sprintf("Path %q for bookmark %q doesn't exist on disk", path, name))
		} else {
			panic(err)
		}
	}

	fmt.Println(path)
}

func (m *Manager) handlePrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  bm --add <name>")
	fmt.Println("  bm --get <name>")
	fmt.Println("  bm <name>")
	os.Exit(0)
}
