package components

type Attack struct {
	Damage int `json:"damage"`
}

func NewAttack(damage int) *Attack {
	return &Attack{Damage: damage}
}

func (a *Attack) Mask() uint64 {
	return MaskAttack
}
