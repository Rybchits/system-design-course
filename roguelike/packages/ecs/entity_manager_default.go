package ecs

type defaultEntityManager struct {
	entities []*Entity
}

func (m *defaultEntityManager) Add(entities ...*Entity) {
	m.entities = append(m.entities, entities...)
}

func (m *defaultEntityManager) Entities() (entities []*Entity) {
	return m.entities
}

func (m *defaultEntityManager) FilterByMask(mask uint64) (entities []*Entity) {

	entities = make([]*Entity, len(m.entities))
	index := 0
	for _, e := range m.entities {

		observed := e.Mask()

		if observed&mask == mask {

			entities[index] = e
			index++
		}
	}

	return entities[:index]
}

func (m *defaultEntityManager) FilterByNames(names ...string) (entities []*Entity) {

	entities = make([]*Entity, len(m.entities))
	index := 0
	for _, e := range m.entities {

		matched := 0
		for _, name := range names {
			for _, c := range e.Components {
				switch v := c.(type) {
				case ComponentWithName:
					if v.Name() == name {
						matched++
					}
				}
			}
		}

		if matched == len(names) {

			entities[index] = e
			index++
		}
	}

	return entities[:index]
}

func (m *defaultEntityManager) Get(id string) (entity *Entity) {
	for _, e := range m.entities {
		if e.Id == id {
			return e
		}
	}
	return
}

func (m *defaultEntityManager) Remove(entity *Entity) {
	for i, e := range m.entities {
		if e.Id == entity.Id {
			copy(m.entities[i:], m.entities[i+1:])
			m.entities[len(m.entities)-1] = nil
			m.entities = m.entities[:len(m.entities)-1]
			break
		}
	}
}

func NewEntityManager() *defaultEntityManager {
	return &defaultEntityManager{
		entities: []*Entity{},
	}
}
