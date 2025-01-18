package components

// Компонент фракции, необходим для отличия противников от союзников (используется в боевой системе)
type Fraction uint64

const (
	FriendsFraction = 0
	EnemiesFraction = 1
)

func (t *Fraction) Mask() uint64 {
	return MaskFraction
}

func NewFraction(value int) *Fraction {
	t := Fraction(value)
	return &t
}
