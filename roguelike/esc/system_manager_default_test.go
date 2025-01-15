package ecs_test

import (
	"testing"

	ecs "roguelike/esc"
)

func TestSystemManager_Systems_Should_Have_No_System_At_Start(t *testing.T) {
	m := ecs.NewSystemManager()
	if len(m.Systems()) != 0 {
		t.Errorf("SystemManager should have no system at start, but got %d", len(m.Systems()))
	}
}

func TestSystemManager_Systems_Should_Have_One_System_After_Adding_One_System(t *testing.T) {
	m := ecs.NewSystemManager()
	s := &mockupDedicatedSystem{}
	m.Add(s)
	if len(m.Systems()) != 1 {
		t.Errorf("SystemManager should have one system at start, but got %d", len(m.Systems()))
	}
}

func TestSystemManager_Systems_Should_Have_Two_System_After_Adding_Two_System(t *testing.T) {
	m := ecs.NewSystemManager()
	s1 := &mockupDedicatedSystem{}
	s2 := &mockupDedicatedSystem{}
	m.Add(s1, s2)
	if len(m.Systems()) != 2 {
		t.Errorf("SystemManager should have one system at start, but got %d", len(m.Systems()))
	}
}

type mockupDedicatedSystem struct{}

func (s *mockupDedicatedSystem) Process(entityManager ecs.EntityManager) (state int) {
	return ecs.StateEngineContinue
}
func (s *mockupDedicatedSystem) Setup()    {}
func (s *mockupDedicatedSystem) Teardown() {}
