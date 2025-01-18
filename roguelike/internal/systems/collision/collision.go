package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type CollisionHandler interface {
	CanHandle(entity1, entity2 *ecs.Entity) bool

	// возвращает true если нужно ли продолжать игру
	Handle(entity1, entity2 *ecs.Entity) bool
}

// Система определения пересечения двух сущностей на карте
type collisionSystem struct {
	handlers []CollisionHandler
}

func (a *collisionSystem) Teardown() {}

func (a *collisionSystem) Setup() {}

func (s *collisionSystem) Process(em ecs.EntityManager) (state int) {
	// Получаем все сущности с компонентом движения
	movingEntities := em.FilterByMask(components.MaskMovement)

	// Получаем все сущности с компонентом позиции
	entities := em.FilterByMask(components.MaskPosition)

	for _, entity1 := range movingEntities {
		movement := entity1.Get(components.MaskMovement).(*components.Movement)

		for _, entity2 := range entities {
			position2 := entity2.Get(components.MaskPosition).(*components.Position)

			// Если это пересечение с тем же самым персонажем или пересечения нет - пропускаем
			if entity1 == entity2 || movement.Next.X != position2.X || movement.Next.Y != position2.Y {
				continue
			}

			// Запускаем обработчики столкновений
			for _, handler := range s.handlers {
				if handler.CanHandle(entity1, entity2) {
					engineContinue := handler.Handle(entity1, entity2)

					// Две сущности не могут находится на одной позиции, поэтому отменяем перемещение
					entity1.Remove(components.MaskMovement)

					if !engineContinue {
						return ecs.StateEngineStop
					}
					break
				}
			}
		}
	}

	return ecs.StateEngineContinue
}

func (a *collisionSystem) WithHandlers(handlers ...CollisionHandler) *collisionSystem {
	a.handlers = handlers
	return a
}

func NewCollisionSystem() *collisionSystem {
	return &collisionSystem{}
}
