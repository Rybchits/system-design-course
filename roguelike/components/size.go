package components

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (a *Size) Mask() uint64 {
	return MaskSize
}

func (a *Size) WithWidth(width int) *Size {
	a.Width = width
	return a
}

func (a *Size) WithHeight(height int) *Size {
	a.Height = height
	return a
}

func NewSize() *Size {
	return &Size{}
}
