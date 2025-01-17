package components

type Experience struct {
	CurrentXP int
	Level     int
}

func NewExperience() *Experience {
	return &Experience{
		CurrentXP: 0,
		Level:     1,
	}
}

func (a *Experience) Mask() uint64 {
	return MaskExperience
}

func (e *Experience) AddXP(amount int) {
	e.CurrentXP += amount
}

func (e *Experience) LevelUp(requiredExperience int) {
	e.Level++
	e.CurrentXP -= requiredExperience
}
