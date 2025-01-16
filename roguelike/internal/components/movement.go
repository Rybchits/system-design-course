package components

type Movement struct {
	Previous Position `json:"prev"`
	Next     Position `json:"next"`
}

func (a *Movement) Mask() uint64 {
	return MaskMovement
}

func (a *Movement) WithPrevious(previous Position) *Movement {
	a.Previous = previous
	return a
}

func (a *Movement) WithNext(next Position) *Movement {
	a.Next = next
	return a
}

func NewMovement() *Movement {
	return &Movement{}
}
