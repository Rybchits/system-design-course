package ai

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type mobsBehaviorSystem struct{}

func (s *mobsBehaviorSystem) Process(em ecs.EntityManager) (state int) {
	movementEntities := em.FilterByMask(components.MaskMobStrategy)

	for _, entity := range movementEntities {
		strategy := entity.Get(components.MaskMobStrategy).(*components.MobBehavior).Strategy
		strategy.Act(entity, em)
	}
	return ecs.StateEngineContinue
}

func (a *mobsBehaviorSystem) Teardown() {}

func (a *mobsBehaviorSystem) Setup() {}

func NewmobsBehaviorSystem() *mobsBehaviorSystem {
	return &mobsBehaviorSystem{}
}
