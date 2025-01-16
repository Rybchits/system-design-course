package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"roguelike/internal/components"
	"roguelike/internal/systems"
	inputSystemPackage "roguelike/internal/systems/input"
	ecs "roguelike/packages/ecs"

	"github.com/gdamore/tcell/v2"
)

type defaultGameBuilder struct {
	location components.Location
	screen   tcell.Screen
	engine   ecs.Engine
}

func NewDefaultGameBuilder() *defaultGameBuilder {
	return &defaultGameBuilder{}
}

func (b *defaultGameBuilder) SetLocation(location string) {
	filePath := fmt.Sprintf("resources/location_%s.json", location)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read location file: %v", err)
		os.Exit(1)
	}
	var locationData components.Location
	if err := json.Unmarshal(file, &locationData); err != nil {
		log.Fatalf("Failed to unmarshal location data: %v", err)
		os.Exit(1)
	}
	b.location = locationData
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
	width, height := b.location.MapSize.Width, b.location.MapSize.Height
	sm := ecs.NewSystemManager()
	em := ecs.NewEntityManager()

	sm.Add(
		systems.NewRenderingSystem().WithWidth(width).WithHeight(height).WithScreen(&b.screen),
		inputSystemPackage.NewInputSystem().WithScreen(&b.screen).WithInputHandlers(
			inputSystemPackage.NewMoveHandler(),
			inputSystemPackage.NewQuitHandler(),
		),
	)

	for index, obstacle := range b.location.Obstacles {
		id := fmt.Sprintf("obstacle-%d", index)
		obstacleEntity := ecs.NewEntity(id, []ecs.Component{
			components.NewPosition().WithX(obstacle.Pos.X).WithY(obstacle.Pos.Y),
			components.NewTexture('#'),
		})
		em.Add(obstacleEntity)
	}

	player := ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(b.location.StartPosition.X).WithY(b.location.StartPosition.Y),
		components.NewTexture('@'),
	})
	em.Add(player)

	location := ecs.NewEntity("location", []ecs.Component{&b.location})
	em.Add(location)

	b.engine = ecs.NewDefaultEngine(em, sm)
	b.engine.Setup()
}

func (b *defaultGameBuilder) GetResult() GameModel {
	return defaultGameModel{
		screen: b.screen,
		engine: b.engine,
	}
}
