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

    if test (count $argv) -eq 0 || test $argv[1] = "-h" || test $argv[1] = "--help"
        set -l help_output (bookmarker --help)
        string replace bookmarker bm $help_output
        return 0
    end

    bookmarker $argv
end
