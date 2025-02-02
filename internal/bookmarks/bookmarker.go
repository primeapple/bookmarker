package bookmarker

type Manager struct {
    config *config.Config
    store  Storage
}

func NewManager(cfg *config.Config) *Manager {
    return &Manager{
        config: cfg,
        store:  NewJSONStorage(cfg.StoragePath),
    }
}

func (m *Manager) Run(args []string) error {
    if len(args) == 0 {
        return m.printUsage()
    }

    switch args[0] {
    case "--add":
        return m.handleAdd(args[1:])
    case "--list":
        return m.handleList(args[1:])
    case "--remove":
        return m.handleRemove(args[1:])
    default:
        return m.handleGoto(args[0])
    }
}
