package level

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

func OnDamageCallback(attacker, attacked *ecs.Entity) {
	if attacker.Get(components.MaskExperience) == nil || attacked.Get(components.MaskHealth) == nil {
		return
	}
	experience := attacker.Get(components.MaskExperience).(*components.Experience)
	health := attacked.Get(components.MaskHealth).(*components.Health)

	if !health.IsAlive() {
		experience.AddXP(5)
	}
}
