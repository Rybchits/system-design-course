package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type combatSystem struct{}

func (a *combatSystem) Teardown() {}

func (a *combatSystem) Setup() {}

func (a *combatSystem) Process(em ecs.EntityManager) (state int) {
	movementEntities := em.FilterByMask(components.MaskMovement)
	entities := em.FilterByMask(components.MaskPosition)

	for _, movementEntity := range movementEntities {
		movement := movementEntity.Get(components.MaskMovement).(*components.Movement)

		for _, entity := range entities {
			entityPosition := entity.Get(components.MaskPosition).(*components.Position)

			// TODO переделать на нанесение урона
			// Если есть пересечение позиций, то отменяем переход
			if entityPosition.X == movement.Next.X && entityPosition.Y == movement.Next.Y {
				movementEntity.Remove(components.MaskMovement)
				break
			}
		}
	}
	return ecs.StateEngineContinue
}

func NewCombatSystem() *combatSystem {
	return &combatSystem{}
}
