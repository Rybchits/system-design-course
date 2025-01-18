package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
)

// Обработчик столкновения ведущий к бою
type AttackHandler struct {
	OnDamageCallback func(attacker, attacked *ecs.Entity)
}

func NewAttackHandler() *AttackHandler {
	return &AttackHandler{}
}

func (h *AttackHandler) WithOnDamageCallback(callback func(attacker, attacked *ecs.Entity)) *AttackHandler {
	h.OnDamageCallback = callback
	return h
}

// Можно атаковать, если наступающий имеет атаку, а на которого наступают имеет здоровье
func (h *AttackHandler) CanHandle(entity1, entity2 *ecs.Entity) bool {
	return entity1.Get(components.MaskAttack) != nil && entity2.Get(components.MaskHealth) != nil
}

func haveSameFraction(entity1, entity2 *ecs.Entity) bool {
	fraction1 := entity1.Get(components.MaskFraction).(*components.Fraction)
	fraction2 := entity2.Get(components.MaskFraction).(*components.Fraction)
	return fraction1 == fraction2
}

func (h *AttackHandler) Handle(entity1, entity2 *ecs.Entity) bool {
	if haveSameFraction(entity1, entity2) {
		return true
	}

	health1 := entity1.Get(components.MaskHealth).(*components.Health)
	attack1 := entity1.Get(components.MaskAttack).(*components.Attack)

	health2 := entity2.Get(components.MaskHealth).(*components.Health)
	attack2 := entity2.Get(components.MaskAttack).(*components.Attack)

	// Отнимаем здоровье у сущности на которую напали
	health2.TakeDamage(attack1.Damage)
	h.OnDamageCallback(entity1, entity2)

	// Отнимает здоровье у нападавшего
	if health1 != nil && attack1 != nil {
		health1.TakeDamage(attack2.Damage)
		h.OnDamageCallback(entity2, entity1)
	}

	// Удаляем нападающего с карты, если она была убита
	if !health1.IsAlive() {
		entity1.Remove(components.MaskPosition)
	}

	// Удаляем сущности на которую напали с карты, если она была убита
	if !health2.IsAlive() {
		entity2.Remove(components.MaskPosition)
	}

	// Если сущность игрока была убита, останавливаем игру
	isPlayerDead := (entity1.Id == "player" && !health1.IsAlive()) || (entity2.Id == "player" && !health2.IsAlive())
	return !isPlayerDead
}
