package components

type Texture rune

func (t *Texture) Mask() uint64 {
	return MaskTexture
}

func NewTexture(a rune) *Texture {
	t := Texture(a)
	return &t
}
