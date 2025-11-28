set -g __bm_commands add get go help list remove

function bm
    if test (count $argv) -eq 2 && test $argv[1] = go
        set -l path (command bm get $argv[2])
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

complete -c bm --no-file
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments add --description "Add named bookmark to current directory"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments get --description "Print path for named bookmark"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments go --description "Change directory to bookmark path"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments help --description "Displays help for the 'bm' command"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments list --description "List all available bookmarks with their index and path"
complete -c bm --condition "not __fish_seen_subcommand_from $__bm_commands" --arguments remove --description "Remove named bookmark"
complete -c bm --condition "__fish_seen_subcommand_from add" --force-files
complete -c bm --condition "__fish_seen_subcommand_from get" --arguments "(bm list | cut -d'|' -f2 | tr -d ' ')"
complete -c bm --condition "__fish_seen_subcommand_from go" --arguments "(bm list | cut -d'|' -f2 | tr -d ' ')"
complete -c bm --condition "__fish_seen_subcommand_from remove" --arguments "(bm list | cut -d'|' -f2 | tr -d ' ')"
