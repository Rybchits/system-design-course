package entities

import ecs "roguelike/packages/ecs"

type EntityFactory interface {
	CreatePlayer(x, y int, health int, attack int) *ecs.Entity
	CreateEnemy(entityId string, typeEnemy string, x, y int, health int, attack int) *ecs.Entity
}
