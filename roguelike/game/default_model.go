package game

import (
	ecs "roguelike/esc"

	"github.com/gdamore/tcell/v2"
)

type defaultGameModel struct {
	engine ecs.Engine
	screen tcell.Screen
}

func (g defaultGameModel) Run() {
	g.engine.Run()
}

func (g defaultGameModel) Stop() {
	g.screen.Fini()
	g.engine.Teardown()
}
