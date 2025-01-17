package components

import "roguelike/packages/ecs"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (a *Position) Mask() uint64 {
	return MaskPosition
}

func (a *Position) WithX(x int) *Position {
	a.X = x
	return a
}

func (a *Position) WithY(y int) *Position {
	a.Y = y
	return a
}

func (a *Position) IsFree(em ecs.EntityManager) bool {
	entities := em.FilterByMask(MaskPosition)

	for _, entity := range entities {
		position := entity.Get(MaskPosition).(*Position)
		if position.X == a.X && position.Y == a.Y {
			return false
		}
	}
	return true
}

func (a *Position) Clone() *Position {
	return (&Position{}).WithX(a.X).WithY(a.Y)
}

func NewPosition() *Position {
	return &Position{}
}
