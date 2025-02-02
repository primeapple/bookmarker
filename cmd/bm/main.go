package main

import (
    "fmt"
    "os"

    "github.com/yourusername/directory-bookmarks/internal/bookmarks"
    "github.com/yourusername/directory-bookmarks/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
        os.Exit(1)
    }

    app := bookmarks.NewManager(cfg)
    if err := app.Run(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
