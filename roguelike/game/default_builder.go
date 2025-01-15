package game

import (
	"log"
	"os"
	"roguelike/components"
	ecs "roguelike/esc"
	"roguelike/systems"

	"github.com/gdamore/tcell/v2"
)

type defaultGameBuilder struct {
	level  int
	screen tcell.Screen
	engine ecs.Engine
}

func NewDefaultGameBuilder() *defaultGameBuilder {
	return &defaultGameBuilder{}
}

func (b *defaultGameBuilder) SetLevel(level int) {
	b.level = level
}

func (b *defaultGameBuilder) BuildScreen() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
		os.Exit(1)
	}
	b.screen = screen
}

func (b *defaultGameBuilder) BuildEngine() {
	width, height := 10, 10
	sm := ecs.NewSystemManager()
	em := ecs.NewEntityManager()

	sm.Add(
		systems.NewRenderingSystem().WithWidth(width).WithHeight(height).WithScreen(&b.screen),
		systems.NewInputSystem().WithScreen(&b.screen),
	)

	em.Add(ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(0).WithY(0),
		components.NewTexture('@'),
	}))

	b.engine = ecs.NewDefaultEngine(em, sm)
	b.engine.Setup()
}

func (b *defaultGameBuilder) GetResult() GameModel {
	return defaultGameModel{
		screen: b.screen,
		engine: b.engine,
	}
}
