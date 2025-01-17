package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
	"testing"
)

type MockCollisionHandler struct {
	canHandle bool
	handle    bool
}

func (m *MockCollisionHandler) CanHandle(entity1, entity2 *ecs.Entity) bool {
	return m.canHandle
}

func (m *MockCollisionHandler) Handle(entity1, entity2 *ecs.Entity) bool {
	return m.handle
}

func TestCollisionSystem_NoCollision(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCollisionSystem()

	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(3).WithY(3),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	position := entity1.Get(components.MaskPosition).(*components.Position)
	if position.X != 1 || position.Y != 1 {
		t.Errorf("Ожидается (1, 1), но получено (%d, %d)", position.X, position.Y)
	}

	if entity1.Get(components.MaskMovement) == nil {
		t.Errorf("Movement компонент не должен быть удален")
	}
}

func TestCollisionSystem_WithCollision(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCollisionSystem().WithHandlers(&MockCollisionHandler{canHandle: true, handle: true})

	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	position := entity1.Get(components.MaskPosition).(*components.Position)
	if position.X != 1 || position.Y != 1 {
		t.Errorf("Ожидается (1, 1), но получено (%d, %d)", position.X, position.Y)
	}

	if entity1.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален")
	}
}

func TestCollisionSystem_StopEngine(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCollisionSystem().WithHandlers(&MockCollisionHandler{canHandle: true, handle: false})

	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
	})
	em.Add(entity1)
	em.Add(entity2)

	state := system.Process(em)

	if state != ecs.StateEngineStop {
		t.Errorf("Ожидается состояние остановки движка, но получено %d", state)
	}
}

func TestCollisionSystem_MultipleHandlers(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCollisionSystem().WithHandlers(
		&MockCollisionHandler{canHandle: false, handle: true},
		&MockCollisionHandler{canHandle: true, handle: true},
	)

	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	position := entity1.Get(components.MaskPosition).(*components.Position)
	if position.X != 1 || position.Y != 1 {
		t.Errorf("Ожидается (1, 1), но получено (%d, %d)", position.X, position.Y)
	}

	if entity1.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален")
	}
}
