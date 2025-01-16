package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type combatSystem struct{}

func (a *combatSystem) Teardown() {}

func (a *combatSystem) Setup() {}

func (a *combatSystem) Process(em ecs.EntityManager) (state int) {
	// Получаем все сущности с компонентом движения
	movementEntities := em.FilterByMask(components.MaskMovement)

	// Получаем все сущности с компонентом позиции
	entities := em.FilterByMask(components.MaskPosition)

	for _, entity1 := range movementEntities {
		movement := entity1.Get(components.MaskMovement).(*components.Movement)
		health1 := entity1.Get(components.MaskHealth).(*components.Health)
		attack1 := entity1.Get(components.MaskAttack).(*components.Attack)

		for _, entity2 := range entities {
			if entity1.Id == entity2.Id {
				continue
			}

			position2 := entity2.Get(components.MaskPosition).(*components.Position)
			health2 := entity2.Get(components.MaskHealth).(*components.Health)
			attack2 := entity2.Get(components.MaskAttack).(*components.Attack)

			// Если есть пересечение позиций, атакуем и отменяем переход
			if position2.X == movement.Next.X && position2.Y == movement.Next.Y {
				health2.TakeDamage(attack1.Damage)
				health1.TakeDamage(attack2.Damage)

				if !health1.IsAlive() {
					entity1.Remove(components.MaskPosition)
				}

				if !health2.IsAlive() {
					entity2.Remove(components.MaskPosition)
				}

				entity1.Remove(components.MaskMovement)
				break
			}
		}
	}
	return ecs.StateEngineContinue
}

func NewCombatSystem() *combatSystem {
	return &combatSystem{}
}
