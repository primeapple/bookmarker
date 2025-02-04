package main

import (
    "fmt"
    "os"

    "github.com/primeapple/bookmarker/internal/manager"
)

func main() {
    app := manager.NewManager()
    if err := app.Run(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
