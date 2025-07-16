function bm
    if test (count $argv) -eq 2 && test $argv[1] = go
        set -l path (bookmarker get $argv[2])
        set -l code $status
        if test $code -ne 0
            return $code
        end

        cd $path
        return 0
    end

    if test (count $argv) -eq 0 || test $argv[1] = help || test $argv[1] = init
        set -l help_output (bookmarker help)
        string replace bookmarker bm $help_output
        return 0
    end

    bookmarker $argv
end

set -l commands add get go help list remove

complete -c bm --no-file
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "add" --description "Add named bookmark to current directory"
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "get" --description "Print path for named bookmark"
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "go" --description "Change directory to bookmark path"
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "help" --description "Displays help for the 'bm' command"
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "list" --description "List all available bookmarks with their index and path"
complete -c bm --condition "not __fish_seen_subcommand_from $commands" --arguments "remove" --description "Remove named bookmark"
complete -c bm --condition "__fish_seen_subcommand_from get" --arguments "(bookmarker list | cut -d'|' -f2 | tr -d ' ')"
complete -c bm --condition "__fish_seen_subcommand_from go" --arguments "(bookmarker list | cut -d'|' -f2 | tr -d ' ')"
complete -c bm --condition "__fish_seen_subcommand_from remove" --arguments "(bookmarker list | cut -d'|' -f2 | tr -d ' ')"
