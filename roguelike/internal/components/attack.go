package components

// Компонент атаки, содержащий количество урона, которое сущность может нанести в бою
type Attack struct {
	Damage int `json:"damage"`
}

func NewAttack(damage int) *Attack {
	return &Attack{Damage: damage}
}

func (a *Attack) Mask() uint64 {
	return MaskAttack
}
