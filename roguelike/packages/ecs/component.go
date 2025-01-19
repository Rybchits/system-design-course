package ecs

// Компонент содержат исключительно данные
type Component interface {
	Mask() uint64
}

// Используется для фильтрации сущностей по имени
type ComponentWithName interface {
	Component
	Name() string
}
