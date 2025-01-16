package components

type Health struct {
	CurrentHealth int
	MaxHealth     int
}

func (a *Health) Mask() uint64 {
	return MaskHealth
}

func NewHealth(maxHealth int) *Health {
	return &Health{
		CurrentHealth: maxHealth,
		MaxHealth:     maxHealth,
	}
}

func (h *Health) TakeDamage(amount int) {
	h.CurrentHealth -= amount
	if h.CurrentHealth < 0 {
		h.CurrentHealth = 0
	}
}

func (h *Health) Heal(amount int) {
	h.CurrentHealth += amount
	if h.CurrentHealth > h.MaxHealth {
		h.CurrentHealth = h.MaxHealth
	}
}

func (h *Health) IsAlive() bool {
	return h.CurrentHealth > 0
}
