package main

import (
	"roguelike/game"
)

func main() {
	builder := game.NewDefaultGameBuilder()
	builder.SetLevel(1)
	builder.BuildScreen()
	builder.BuildEngine()

	game := builder.GetResult()
	defer game.Stop()
	game.Run()
}
