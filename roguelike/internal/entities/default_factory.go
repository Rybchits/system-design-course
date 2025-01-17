package entities

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

type defaultEntityFactory struct{}

func NewDefaultEntityFactory() *defaultEntityFactory {
	return &defaultEntityFactory{}
}

func (f *defaultEntityFactory) CreatePlayer(x, y int, health int, attack int) *ecs.Entity {
	healthComponent := components.NewHealth(health)
	attackComponent := components.NewAttack(attack)

	return ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(x).WithY(y),
		components.NewTexture('@'),
		healthComponent,
		attackComponent,
	})
}

func (f *defaultEntityFactory) CreateEnemy(entityId string, typeEnemy string, x, y int, health int, attack int) *ecs.Entity {
	position := components.NewPosition().WithX(x).WithY(y)
	healthComponent := components.NewHealth(health)
	attackComponent := components.NewAttack(attack)

	var texture *components.Texture
	switch typeEnemy {
	case "dumb":
		texture = components.NewTexture('D')
	case "aggressive":
		texture = components.NewTexture('A')
	case "passive":
		texture = components.NewTexture('P')
	default:
		return nil
	}
	// TODO добавить компонент поведения
	return ecs.NewEntity(entityId, []ecs.Component{position, texture, healthComponent, attackComponent})
}
