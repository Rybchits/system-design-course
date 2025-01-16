package ecs

type defaultSystemManager struct {
	systems []System
}

func (m *defaultSystemManager) Add(systems ...System) {
	m.systems = append(m.systems, systems...)
}

func (m *defaultSystemManager) Systems() []System {
	return m.systems
}

func NewSystemManager() SystemManager {
	return &defaultSystemManager{
		systems: []System{},
	}
}
