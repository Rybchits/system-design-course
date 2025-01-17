package entities

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

type EntityFactory interface {
	CreatePlayer(x, y int, health int, attack int) *ecs.Entity
	CreateEnemy(entityId string, description components.Enemy) *ecs.Entity
}
