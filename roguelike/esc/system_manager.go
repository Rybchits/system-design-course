package ecs

type SystemManager interface {
	Add(systems ...System)

	Systems() []System
}
