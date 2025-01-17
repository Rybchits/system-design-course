package main

import (
	"roguelike/internal/game"
)

func main() {
	builder := game.NewDefaultGameBuilder().WithLocation("resources/location_1.json")
	builder.BuildScreen()
	builder.BuildEngine()

	game := builder.GetResult()
	defer game.Stop()
	game.Run()
}
