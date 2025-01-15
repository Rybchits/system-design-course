package components

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

func (a *Position) Clone() *Position {
	return (&Position{}).WithX(a.X).WithY(a.Y)
}

func NewPosition() *Position {
	return &Position{}
}
