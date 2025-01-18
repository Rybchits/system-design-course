package entities

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

// Дефотная реализация абстрактоной фабрики для создания сущностей
type defaultEntityFactory struct{}

func NewDefaultEntityFactory() *defaultEntityFactory {
	return &defaultEntityFactory{}
}

func (f *defaultEntityFactory) CreatePlayer(x, y int, health int, attack int) *ecs.Entity {
	healthComponent := components.NewHealth(health)
	attackComponent := components.NewAttack(attack)
	fractionComponent := components.NewFraction(components.FriendsFraction)

	return ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(x).WithY(y),
		components.NewTexture('@'),
		components.NewExperience(),
		healthComponent,
		attackComponent,
		fractionComponent,
	})
}

func (f *defaultEntityFactory) CreateEnemy(entityId string, description components.Enemy) *ecs.Entity {
	position := components.NewPosition().WithX(description.Pos.X).WithY(description.Pos.Y)
	healthComponent := components.NewHealth(description.Health)
	attackComponent := components.NewAttack(description.Attack)
	fractionComponent := components.NewFraction(components.EnemiesFraction)

	var texture *components.Texture
	var strategy ecs.Component

	switch description.Type {
	case "aggressive":
		texture = components.NewTexture('A')
		strategy = components.NewMobBehavior(components.NewAggressiveStrategy(500))
	case "passive":
		texture = components.NewTexture('P')
		strategy = components.NewMobBehavior(components.NewPassiveStrategy())
	case "cowardly":
		texture = components.NewTexture('C')
		strategy = components.NewMobBehavior(components.NewCowardStrategy(500))
	default:
		return nil
	}

	return ecs.NewEntity(entityId, []ecs.Component{
		position,
		texture,
		healthComponent,
		attackComponent,
		fractionComponent,
		strategy,
	})
}
