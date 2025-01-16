package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type movementSystem struct{}

func (a *movementSystem) Teardown() {}

func (a *movementSystem) Setup() {}

func (a *movementSystem) Process(em ecs.EntityManager) (state int) {
	movementEntities := em.FilterByMask(components.MaskMovement | components.MaskPosition)

	for _, entity := range movementEntities {
		movement := entity.Get(components.MaskMovement).(*components.Movement)
		position := entity.Get(components.MaskPosition).(*components.Position)

		position.X = movement.Next.X
		position.Y = movement.Next.Y
		entity.Remove(components.MaskMovement)
	}

	return ecs.StateEngineContinue
}

func NewMovementSystem() *movementSystem {
	return &movementSystem{}
}
