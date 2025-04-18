package shell

import (
	_ "embed"
)

//go:embed resources/init.fish
var initFish []byte

func InitFish() string {
	return string(initFish)
}
