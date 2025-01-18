package components

// Маски компонентов, необходимые для эффективного получения сущностей
const (
	MaskPosition    = uint64(1 << 0)
	MaskTexture     = uint64(1 << 1)
	MaskLocation    = uint64(1 << 2)
	MaskMovement    = uint64(1 << 3)
	MaskAttack      = uint64(1 << 4)
	MaskHealth      = uint64(1 << 5)
	MaskFraction    = uint64(1 << 6)
	MaskMobStrategy = uint64(1 << 7)
	MaskExperience  = uint64(1 << 8)
)
