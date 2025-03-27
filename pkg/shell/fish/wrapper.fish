function bm
    if test (count $argv) -eq 1 && not string match --quiet --regex -- '-.*' $argv[1]
        set -l path (bookmarker --get $argv[1])
        set -l code $status
        if test $code -ne 0
            return $code
        end

        cd $path
        return 0
    end

    if test (count $argv) -eq 0 || test $argv[1] = -h || test $argv[1] = --help
        set -l help_output (bookmarker --help)
        string replace bookmarker bm $help_output
        return 0
    end

    bookmarker $argv
end

set -l options -a --add -g --get -l --list

complete -c bm --no-file
complete -c bm --condition "not __fish_seen_subcommand_from $options" --require-parameter --arguments "(bm --list | awk '{ print \$1 }')" --description "Change directory to bookmark path"
complete -c bm --old-option a --long-option add --description "Add named bookmark to current directory"
complete -c bm --old-option g --long-option get --require-parameter --arguments "(bm --list | awk '{ print \$1 }')" --description "Print path to bookmark"
complete -c bm --old-option h --long-option help --description "Displays help for the 'bm' command"
complete -c bm --old-option l --long-option list --description "List all available bookmarks with their index and path"
