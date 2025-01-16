package main

import (
	"roguelike/internal/game"
)

func main() {
	builder := game.NewDefaultGameBuilder().WithLocation("1")
	builder.BuildScreen()
	builder.BuildEngine()

	game := builder.GetResult()
	defer game.Stop()
	game.Run()
}
