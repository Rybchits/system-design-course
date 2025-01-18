package game

import (
	"encoding/json"
	"fmt"
	"os"
	"roguelike/internal/components"
	"roguelike/internal/entities"
	collisionSystemPackage "roguelike/internal/systems/collision"
	inputSystemPackage "roguelike/internal/systems/input"
	levelSystemPackage "roguelike/internal/systems/level"
	mobsBehaviorSystemPackage "roguelike/internal/systems/mobs_behavior"
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
		entityFactory: entities.NewDefaultEntityFactory(),
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

func (b *defaultGameBuilder) SetLocation(locationFilePath string) error {
	file, err := os.ReadFile(locationFilePath)
	if err != nil {
		return fmt.Errorf("failed to read location file: %v", err)
	}
	var locationData components.Location
	if err := json.Unmarshal(file, &locationData); err != nil {
		return fmt.Errorf("failed to unmarshal location data: %v", err)
	}
	b.location = locationData
	return nil
}

func (b *defaultGameBuilder) BuildScreen() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to initialize screen: %v", err)
	}
	b.screen = screen
	return nil
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
		mobsBehaviorSystemPackage.NewmobsBehaviorSystem(),
		collisionSystemPackage.NewCollisionSystem().WithHandlers(
			collisionSystemPackage.NewAttackHandler().WithOnDamageCallback(levelSystemPackage.OnDamageCallback),
		),
		movementSystemPackage.NewMovementSystem(),
		levelSystemPackage.NewExperienceSystem(),
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
		if entity := b.entityFactory.CreateEnemy(id, enemy); entity != nil {
			em.Add(entity)
		}
	}
	em.Add(ecs.NewEntity("location", []ecs.Component{&b.location}))

	b.engine = ecs.NewDefaultEngine(em, sm)
	b.engine.Setup()
}

func (b *defaultGameBuilder) GetResult() GameModel {
	return defaultGameModel{
		screen: b.screen,
		engine: b.engine,
	}
}
