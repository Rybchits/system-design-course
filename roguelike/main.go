package main

import (
	"log"
	"os"
	"roguelike/components"
	ecs "roguelike/esc"
	"roguelike/systems"

	"github.com/gdamore/tcell/v2"
)

func main() {
	width, height := 10, 10
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
		os.Exit(1)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
		os.Exit(1)
	}

	sm := ecs.NewSystemManager()
	em := ecs.NewEntityManager()

	sm.Add(
		systems.NewRenderingSystem().WithWidth(width).WithHeight(height).WithScreen(&screen),
		systems.NewInputSystem().WithScreen(&screen),
	)

	em.Add(ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(0).WithY(0),
		components.NewTexture('@'),
	}))

	de := ecs.NewDefaultEngine(em, sm)
	de.Setup()
	defer de.Teardown()
	de.Run()
}
