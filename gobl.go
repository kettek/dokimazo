package main

import (
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	runArgs := append([]interface{}{}, "./dokimazo"+exe)

	Task("build").
		Exec("go", "build", "./cmd/dokimazo")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("cmd/dokimazo/*", "internal/game/*", "internal/res/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Go()
}
