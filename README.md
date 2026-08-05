# bookmarker
It's [harpoon](https://github.com/ThePrimeagen/harpoon) for your shell.

## Installation

1. Install the `bm` binary

### Releases

```sh
go install github.com/primeapple/bookmarker/cmd/bm@latest
```

### Locally

```sh
make install
```

## Setup

Connect your shell
    1. Fish
    ```sh
    # add to your config.fish
    bm init fish | source
    ```

## Usage

A bookmark is a name pointing at an unordered set of directories. Most
bookmarks hold a single directory, but you can group related ones, for example
all repositories of a project, under one name.

| Command                     | Description                                                        |
| --------------------------- | ------------------------------------------------------------------ |
| `bm add <name> <path>...`   | Create a bookmark, replacing all paths it currently holds          |
| `bm attach <name> <path>...`| Add paths to a bookmark, creating it if it doesn't exist           |
| `bm detach <name> <path>...`| Remove paths from a bookmark, deleting it if no path is left       |
| `bm remove <name>...`       | Delete whole bookmarks                                             |
| `bm get <name>`             | Print every path of a bookmark, one per line                       |
| `bm go <name>`              | Change directory to the bookmark's path                            |
| `bm list [--porcelain]`     | List all bookmarks, as a table or as `name<TAB>path` lines         |

```sh
bm add ui ~/src/ui                # ui -> ~/src/ui
bm go ui                          # cd ~/src/ui
bm attach ui ~/src/design-system  # ui -> ~/src/ui ~/src/design-system
bm get ui                         # prints both paths
bm go ui                          # error: a bookmark with 2 paths is no jump target
bm detach ui ~/src/design-system  # ui -> ~/src/ui
bm go ui                          # works again
```

`attach` and `detach` treat the paths as a set: paths are stored absolute,
duplicates are dropped, and order is not preserved.

`bm go` only works on bookmarks holding exactly one path. `bm get` is the
scriptable primitive and never touches the disk:

```sh
for dir in (bm get ui)
    git -C $dir status --short
end
```

