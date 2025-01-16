package entities

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

type EntityFactory struct{}

func NewEntityFactory() *EntityFactory {
	return &EntityFactory{}
}

func (f *EntityFactory) CreatePlayer(x, y int) *ecs.Entity {
	return ecs.NewEntity("player", []ecs.Component{
		components.NewPosition().WithX(x).WithY(y),
		components.NewTexture('@'),
	})
}

func (f *EntityFactory) CreateLocation(location components.Location) *ecs.Entity {
	return ecs.NewEntity("location", []ecs.Component{&location})
}

func (f *EntityFactory) CreateEnemy(entityId string, typeEnemy string, x, y int, health int, attack int) *ecs.Entity {
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
