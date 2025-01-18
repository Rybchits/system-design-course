package level

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
)

// Система для начисления опыта и перехода на новый уровень
type experienceSystem struct{}

func NewExperienceSystem() *experienceSystem {
	return &experienceSystem{}
}

func (s *experienceSystem) Setup() {}

// Проходит по всем сущностям с компонентом Experience и проверяет, достиг ли текущий опыт уровня
func (s *experienceSystem) Process(em ecs.EntityManager) int {
	entities := em.FilterByMask(components.MaskExperience)

	for _, entity := range entities {
		experience := entity.Get(components.MaskExperience).(*components.Experience)

		// Проверка на уровень. Если текущий опыт больше или равен уровню * 10, то увеличиваем уровень на 1
		// и увеличиваем максимальное здоровье на 10
		for experience.CurrentXP >= experience.Level*10 {
			experience.LevelUp(experience.Level * 10)
			if entity.Get(components.MaskHealth) != nil {
				health := entity.Get(components.MaskHealth).(*components.Health)
				health.WithMaxHealth(health.MaxHealth + 10)
			}
		}
	}
	return ecs.StateEngineContinue
}

func (s *experienceSystem) Teardown() {}
