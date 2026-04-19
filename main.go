package main

import (
	"github.com/timkelleher/spaceui/internal/state"
	"github.com/timkelleher/spaceui/internal/ui"
)

func main() {
	state.Init()

	app := ui.NewApp()
	app.Run()
}
