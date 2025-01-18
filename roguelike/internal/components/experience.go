package components

// Текущий опыт персонажа и его уровень
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

// Добавляет опыт персонажу
func (e *Experience) AddXP(amount int) {
	e.CurrentXP += amount
}

// Повышеает уровень персонажа, отнимая у него requiredExperience
func (e *Experience) LevelUp(requiredExperience int) {
	if requiredExperience >= e.CurrentXP {
		e.Level++
		e.CurrentXP -= requiredExperience
	}
}
