set -g __bm_commands add attach detach get go help init list remove

function bm
    if test (count $argv) -eq 2 && test $argv[1] = go
        set -l path (command bm go $argv[2])
        set -l code $status
        if test $code -ne 0
            return $code
        end

        cd $path
        return 0
    end

    command bm $argv
    return $status
end

function __bm_names
    command bm list --porcelain | cut -f1 | uniq
end

# Index of the argument that is currently being completed, ignoring the `bm`
# command itself. `bm attach <here>` is index 1, `bm attach name <here>` is 2.
function __bm_arg_index
    math (count (commandline -opc)) - 1
end

function __bm_completing_name
    test (__bm_arg_index) -eq 1
end

# Paths attached to the bookmark named on the current commandline.
function __bm_attached_paths
    set -l tokens (commandline -opc)
    if test (count $tokens) -lt 3
        return
    end

    command bm list --porcelain | string replace --filter --regex "^"(string escape --style=regex $tokens[3])\t ""
end

complete -c bm --no-file
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments add --description "Add named bookmark, replacing all its paths"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments attach --description "Attach additional paths to a bookmark"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments detach --description "Detach paths from a bookmark"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments get --description "Print all paths of a bookmark"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments go --description "Change directory to bookmark path"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments help --description "Displays help for the 'bm' command"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments init --description "Print the shell integration script"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments list --description "List all bookmarks with their paths"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments remove --description "Remove whole bookmarks"

complete -c bm --condition "__fish_seen_subcommand_from attach detach; and __bm_completing_name" --arguments "(__bm_names)"
complete -c bm --condition "__fish_seen_subcommand_from get go remove" --arguments "(__bm_names)"
complete -c bm --condition "__fish_seen_subcommand_from add attach; and not __bm_completing_name" --force-files
complete -c bm --condition "__fish_seen_subcommand_from detach; and not __bm_completing_name" --arguments "(__bm_attached_paths)"

complete -c bm --condition "__fish_seen_subcommand_from init" --arguments fish --description "Fish shell"
complete -c bm --condition "__fish_seen_subcommand_from list" --long porcelain --description "Print machine readable name<TAB>path lines"
