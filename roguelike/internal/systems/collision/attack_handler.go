package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

type AttackHandler struct{}

func NewAttackHandler() *AttackHandler {
	return &AttackHandler{}
}

// Можно атаковать, если наступающий имеет атаку, а на которого наступают имеет здоровье
func (h *AttackHandler) CanHandle(entity1, entity2 *ecs.Entity) bool {
	return entity1.Get(components.MaskAttack) != nil && entity2.Get(components.MaskHealth) != nil
}

func (h *AttackHandler) Handle(entity1, entity2 *ecs.Entity) bool {
	health1 := entity1.Get(components.MaskHealth).(*components.Health)
	attack1 := entity1.Get(components.MaskAttack).(*components.Attack)

	health2 := entity2.Get(components.MaskHealth).(*components.Health)
	attack2 := entity2.Get(components.MaskAttack).(*components.Attack)

	health2.TakeDamage(attack1.Damage)

	if health1 != nil && attack1 != nil {
		health1.TakeDamage(attack2.Damage)
	}

	if !health1.IsAlive() {
		entity1.Remove(components.MaskPosition)
	}

	if !health2.IsAlive() {
		entity2.Remove(components.MaskPosition)
	}

	isPlayerDead := (entity1.Id == "player" && !health1.IsAlive()) || (entity2.Id == "player" && !health2.IsAlive())
	return !isPlayerDead
}
