package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"roguelike/components"
	ecs "roguelike/esc"
	"roguelike/systems"

	"github.com/gdamore/tcell/v2"
)

type defaultGameBuilder struct {
	levelData components.Level
	screen    tcell.Screen
	engine    ecs.Engine
}

func NewDefaultGameBuilder() *defaultGameBuilder {
	return &defaultGameBuilder{}
}

func (b *defaultGameBuilder) SetLevel(level int) {
	filePath := fmt.Sprintf("resources/level_%d.json", level)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read level file: %v", err)
		os.Exit(1)
	}
	var levelData components.Level
	if err := json.Unmarshal(file, &levelData); err != nil {
		log.Fatalf("Failed to unmarshal level data: %v", err)
		os.Exit(1)
	}
	b.levelData = levelData
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
	width, height := b.levelData.MapSize.Width, b.levelData.MapSize.Height
	sm := ecs.NewSystemManager()
	em := ecs.NewEntityManager()

	sm.Add(
		systems.NewRenderingSystem().WithWidth(width).WithHeight(height).WithScreen(&b.screen),
		systems.NewInputSystem().WithScreen(&b.screen),
	)

	for _, obstacle := range b.levelData.Obstacles {
		obstacleEntity := ecs.NewEntity("obstacle", []ecs.Component{
			components.NewPosition().WithX(obstacle.Pos.X).WithY(obstacle.Pos.Y),
			components.NewTexture('#'),
		})
		em.Add(obstacleEntity)
	}

	player := ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(b.levelData.StartPosition.X).WithY(b.levelData.StartPosition.Y),
		components.NewTexture('@'),
	})
	em.Add(player)

	b.engine = ecs.NewDefaultEngine(em, sm)
	b.engine.Setup()
}

func (b *defaultGameBuilder) GetResult() GameModel {
	return defaultGameModel{
		screen: b.screen,
		engine: b.engine,
	}
}
