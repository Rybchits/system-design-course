package game

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"roguelike/internal/components"
	"roguelike/internal/entities"
	collisionSystemPackage "roguelike/internal/systems/collision"
	inputSystemPackage "roguelike/internal/systems/input"
	movementSystemPackage "roguelike/internal/systems/movement"
	renderingSystemPackage "roguelike/internal/systems/rendering"
	ecs "roguelike/packages/ecs"

	"github.com/gdamore/tcell/v2"
)

type defaultGameBuilder struct {
	location      components.Location
	playerAttack  int
	playerHealth  int
	screen        tcell.Screen
	engine        ecs.Engine
	entityFactory entities.EntityFactory
}

func NewDefaultGameBuilder() *defaultGameBuilder {
	return &defaultGameBuilder{
		entityFactory: *entities.NewEntityFactory(),
		playerAttack:  10,
		playerHealth:  100,
	}
}

func (b *defaultGameBuilder) WithPlayerAttack(attack int) *defaultGameBuilder {
	b.playerAttack = attack
	return b
}

func (b *defaultGameBuilder) WithPlayerHealth(health int) *defaultGameBuilder {
	b.playerHealth = health
	return b
}

func (b *defaultGameBuilder) WithLocation(location string) *defaultGameBuilder {
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
	return b
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
		collisionSystemPackage.NewCollisionSystem().WithHandlers(
			collisionSystemPackage.NewAttackHandler(),
		),
		movementSystemPackage.NewMovementSystem(),
		renderingSystemPackage.NewRenderingSystem().WithScreen(&b.screen),
	)

	em.Add(b.entityFactory.CreatePlayer(
		b.location.StartPosition.X,
		b.location.StartPosition.Y,
		b.playerHealth,
		b.playerAttack,
	))

	// Заполняет карту противниками
	for index, enemy := range b.location.Enemies {
		id := fmt.Sprintf("enemy-%d", index)
		entity := b.entityFactory.CreateEnemy(id, enemy.Type, enemy.Pos.X, enemy.Pos.Y, enemy.Health, enemy.Attack)
		if entity != nil {
			em.Add(entity)
		}
	}
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
