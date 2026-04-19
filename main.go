package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/timkelleher/spaceui/internal/state"
	"github.com/timkelleher/spaceui/internal/ui"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	state.Init()

	app := ui.NewApp()
	app.Run()
}
