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

	runArgs := append([]interface{}{}, "./rotate"+exe)

	Task("build").
		Exec("go", "build", "./cmd/rotate")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("cmd/rotate/*", "internal/game/*", "internal/res/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Go()
}
