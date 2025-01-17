package components

type Movement struct {
	Next Position `json:"next"`
}

func (a *Movement) Mask() uint64 {
	return MaskMovement
}

func (a *Movement) WithNext(next Position) *Movement {
	a.Next = next
	return a
}

func NewMovement() *Movement {
	return &Movement{}
}
