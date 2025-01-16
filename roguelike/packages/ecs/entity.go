package ecs

type Entity struct {
	Components []Component `json:"components"`
	Id         string      `json:"id"`
	Masked     uint64      `json:"masked"`
}

// Добавляет компоненты к сущности
func (e *Entity) Add(cn ...Component) {
	for _, c := range cn {
		if e.Masked&c.Mask() == c.Mask() {
			continue
		}
		e.Components = append(e.Components, c)
		e.Masked = maskSlice(e.Components)
	}
}

// Возвращает компонент из сущности
func (e *Entity) Get(mask uint64) Component {
	for _, c := range e.Components {
		if c.Mask() == mask {
			return c
		}
	}
	return nil
}

func (e *Entity) Mask() uint64 {
	return e.Masked
}

// Удаляет компонент из сущности по его маске
func (e *Entity) Remove(mask uint64) {
	modified := false
	for i, c := range e.Components {
		if c.Mask() == mask {
			copy(e.Components[i:], e.Components[i+1:])
			e.Components[len(e.Components)-1] = nil
			e.Components = e.Components[:len(e.Components)-1]
			modified = true
			break
		}
	}
	if modified {
		e.Masked = maskSlice(e.Components)
	}
}

func NewEntity(id string, components []Component) *Entity {
	return &Entity{
		Components: components,
		Id:         id,
		Masked:     maskSlice(components),
	}
}

// Считает маску для сущности по ее компонентам
func maskSlice(components []Component) uint64 {
	mask := uint64(0)
	for _, c := range components {
		mask = mask | c.Mask()
	}
	return mask
}
