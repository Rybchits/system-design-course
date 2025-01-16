package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
	"testing"
)

func TestMovementSystem(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewMovementSystem()

	entity := ecs.NewEntity("test_entity", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	em.Add(entity)

	system.Process(em)

	position := entity.Get(components.MaskPosition).(*components.Position)
	if position.X != 2 || position.Y != 2 {
		t.Errorf("Ожидается (2, 2), но получено (%d, %d)", position.X, position.Y)
	}

	if entity.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален")
	}
}

func TestMovementWithoutMovementComponent(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewMovementSystem()

	entity := ecs.NewEntity("test_entity", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
	})
	em.Add(entity)

	system.Process(em)

	position := entity.Get(components.MaskPosition).(*components.Position)
	if position.X != 1 || position.Y != 1 {
		t.Errorf("Ожидается (1, 1), но получено (%d, %d)", position.X, position.Y)
	}
}

func TestMovementWithMultipleEntities(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewMovementSystem()

	entity1 := ecs.NewEntity("test_entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
	})
	entity2 := ecs.NewEntity("test_entity2", []ecs.Component{
		components.NewPosition().WithX(3).WithY(3),
		components.NewMovement().WithNext(*components.NewPosition().WithX(4).WithY(4)),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	position1 := entity1.Get(components.MaskPosition).(*components.Position)
	if position1.X != 2 || position1.Y != 2 {
		t.Errorf("Ожидается (2, 2) для entity1, но получено (%d, %d)", position1.X, position1.Y)
	}

	position2 := entity2.Get(components.MaskPosition).(*components.Position)
	if position2.X != 4 || position2.Y != 4 {
		t.Errorf("Ожидается (4, 4) для entity2, но получено (%d, %d)", position2.X, position2.Y)
	}

	if entity1.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален для entity1")
	}

	if entity2.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален для entity2")
	}
}
