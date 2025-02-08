function bm
    if test (count $argv) -eq 1 && not string match -- '--'
        set -l path (bookmarker --get $argv[1])
        set -l code $status
        if test status -eq 0
            cd $path
            return 0
        end
        return $code
    end

    bookmarker $argv
end
