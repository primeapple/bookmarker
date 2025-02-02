package storage

type Storage interface {
    Load() (*Bookmarks, error)
    Save(*Bookmarks) error
}

type Bookmarks struct {
    Directories map[string]map[string]string `json:"directories"`
}
