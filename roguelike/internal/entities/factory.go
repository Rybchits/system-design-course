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

func (f *EntityFactory) CreateObstacle(entityId string, typeObstacle string, x, y int) *ecs.Entity {
	position := components.NewPosition().WithX(x).WithY(y)

	var texture *components.Texture
	switch typeObstacle {
	case "wall":
		texture = components.NewTexture('#')
	default:
		return nil
	}
	return ecs.NewEntity(entityId, []ecs.Component{position, texture})
}

func (f *EntityFactory) CreateLocation(location components.Location) *ecs.Entity {
	return ecs.NewEntity("location", []ecs.Component{&location})
}
