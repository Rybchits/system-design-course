package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"roguelike/internal/components"
	"roguelike/internal/entities"
	combatSystemPackage "roguelike/internal/systems/combat"
	inputSystemPackage "roguelike/internal/systems/input"
	movementSystemPackage "roguelike/internal/systems/movement"
	renderingSystemPackage "roguelike/internal/systems/rendering"
	ecs "roguelike/packages/ecs"

	"github.com/gdamore/tcell/v2"
)

type defaultGameBuilder struct {
	location      components.Location
	screen        tcell.Screen
	engine        ecs.Engine
	entityFactory entities.EntityFactory
}

func NewDefaultGameBuilder() *defaultGameBuilder {
	return &defaultGameBuilder{
		entityFactory: *entities.NewEntityFactory(),
	}
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
	sm := ecs.NewSystemManager()
	em := ecs.NewEntityManager()

	// Добавляет системы в менеджер систем
	sm.Add(
		inputSystemPackage.NewInputSystem().WithScreen(&b.screen).WithInputHandlers(
			inputSystemPackage.NewMoveHandler(),
			inputSystemPackage.NewQuitHandler(),
		),
		combatSystemPackage.NewCombatSystem(),
		movementSystemPackage.NewMovementSystem(),
		renderingSystemPackage.NewRenderingSystem().WithScreen(&b.screen),
	)

	// Заполняет карту противниками
	for index, enemy := range b.location.Enemies {
		id := fmt.Sprintf("enemy-%d", index)
		entity := b.entityFactory.CreateEnemy(id, enemy.Type, enemy.Pos.X, enemy.Pos.Y, enemy.Health, enemy.Attack)
		if entity != nil {
			em.Add(entity)
		}
	}
	em.Add(b.entityFactory.CreatePlayer(b.location.StartPosition.X, b.location.StartPosition.Y))
	em.Add(b.entityFactory.CreateLocation(b.location))

	b.engine = ecs.NewDefaultEngine(em, sm)
	b.engine.Setup()
}

func (b *defaultGameBuilder) GetResult() GameModel {
	return defaultGameModel{
		screen: b.screen,
		engine: b.engine,
	}
}
