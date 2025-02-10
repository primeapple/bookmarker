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

    bookmarker $argv
end
