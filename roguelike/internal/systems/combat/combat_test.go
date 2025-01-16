package systems

import (
	"roguelike/internal/components"
	"roguelike/packages/ecs"
	"testing"
)

func TestCombatSystem(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCombatSystem()

	// Создаем две сущности, которые будут атаковать друг друга
	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
		components.NewHealth(100),
		components.NewAttack(10),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
		components.NewHealth(100),
		components.NewAttack(20),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	// Проверяем здоровье после атаки
	health1 := entity1.Get(components.MaskHealth).(*components.Health)
	health2 := entity2.Get(components.MaskHealth).(*components.Health)

	if health1.CurrentHealth != 80 {
		t.Errorf("Ожидается здоровье 80 для entity1, но получено %d", health1.CurrentHealth)
	}

	if health2.CurrentHealth != 90 {
		t.Errorf("Ожидается здоровье 90 для entity2, но получено %d", health2.CurrentHealth)
	}

	// Проверяем, что компонент движения был удален
	if entity1.Get(components.MaskMovement) != nil {
		t.Errorf("Movement компонент должен быть удален для entity1")
	}
}

func TestCombatSystem_EntityDies(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCombatSystem()

	// Создаем две сущности, одна из которых умрет после атаки
	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(2).WithY(2)),
		components.NewHealth(10),
		components.NewAttack(10),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
		components.NewHealth(5),
		components.NewAttack(20),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	// Проверяем здоровье после атаки
	health1 := entity1.Get(components.MaskHealth).(*components.Health)
	health2 := entity2.Get(components.MaskHealth).(*components.Health)

	if health1.CurrentHealth != 0 {
		t.Errorf("Ожидается здоровье 0 для entity1, но получено %d", health1.CurrentHealth)
	}

	if health2.CurrentHealth != 0 {
		t.Errorf("Ожидается здоровье 0 для entity2, но получено %d", health2.CurrentHealth)
	}

	// Проверяем, что сущности были удалены
	if entity1.Get(components.MaskPosition) != nil {
		t.Errorf("Position компонент должен быть удален для entity1")
	}

	if entity2.Get(components.MaskPosition) != nil {
		t.Errorf("Position компонент должен быть удален для entity2")
	}
}

func TestCombatSystem_NoAttack(t *testing.T) {
	em := ecs.NewEntityManager()
	system := NewCombatSystem()

	// Создаем две сущности, которые не будут атаковать друг друга
	entity1 := ecs.NewEntity("entity1", []ecs.Component{
		components.NewPosition().WithX(1).WithY(1),
		components.NewMovement().WithNext(*components.NewPosition().WithX(3).WithY(3)),
		components.NewHealth(100),
		components.NewAttack(10),
	})
	entity2 := ecs.NewEntity("entity2", []ecs.Component{
		components.NewPosition().WithX(2).WithY(2),
		components.NewHealth(100),
		components.NewAttack(20),
	})
	em.Add(entity1)
	em.Add(entity2)

	system.Process(em)

	// Проверяем здоровье после атаки
	health1 := entity1.Get(components.MaskHealth).(*components.Health)
	health2 := entity2.Get(components.MaskHealth).(*components.Health)

	if health1.CurrentHealth != 100 {
		t.Errorf("Ожидается здоровье 100 для entity1, но получено %d", health1.CurrentHealth)
	}

	if health2.CurrentHealth != 100 {
		t.Errorf("Ожидается здоровье 100 для entity2, но получено %d", health2.CurrentHealth)
	}

	// Проверяем, что компонент движения не был удален
	if entity1.Get(components.MaskMovement) == nil {
		t.Errorf("Movement компонент не должен быть удален для entity1")
	}
}
