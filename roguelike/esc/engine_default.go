package ecs

// это просто композиция defaultEntityManager и defaultSystemManager.
type defaultEngine struct {
	entityManager EntityManager
	systemManager SystemManager
}

func (e *defaultEngine) Run() {
	shouldStop := false
	for !shouldStop {
		for _, system := range e.systemManager.Systems() {
			state := system.Process(e.entityManager)
			if state == StateEngineStop {
				shouldStop = true
				break
			}
		}
	}
}

func (e *defaultEngine) Tick() {
	for _, system := range e.systemManager.Systems() {
		if state := system.Process(e.entityManager); state == StateEngineStop {
			break
		}
	}
}

func (e *defaultEngine) Setup() {
	for _, sys := range e.systemManager.Systems() {
		sys.Setup()
	}
}

func (e *defaultEngine) Teardown() {
	for _, sys := range e.systemManager.Systems() {
		sys.Teardown()
	}
}

func NewDefaultEngine(entityManager EntityManager, systemManager SystemManager) Engine {
	return &defaultEngine{
		entityManager: entityManager,
		systemManager: systemManager,
	}
}
