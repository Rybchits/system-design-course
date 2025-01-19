package systems

import (
	"fmt"
	"math"
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Отрисовывает карту и сущности на ней
type renderingSystem struct {
	// Сцена для отрисовки содержимого
	screen *tcell.Screen
}

func (a *renderingSystem) Setup() {}

func (a *renderingSystem) Process(em ecs.EntityManager) (state int) {
	(*a.screen).Clear()
	location := em.Get("location").Get(components.MaskLocation).(*components.Location)

	titleLines, titleLength := a.renderLevelTitle(em, 0, 0)
	offsetX, offsetY := 0, titleLines

	a.renderMap(em, offsetX, offsetY, location.MapSize)
	offsetX += int(math.Max(float64(location.MapSize.Width+1), float64(titleLength+1)))

	a.renderEntityDescription(em, offsetX, offsetY)

	(*a.screen).Show()
	time.Sleep(50 * time.Millisecond)
	return ecs.StateEngineContinue
}

func (a *renderingSystem) renderLevelTitle(em ecs.EntityManager, offsetX, offsetY int) (int, int) {
	experience := em.Get("player").Get(components.MaskExperience).(*components.Experience)
	title := fmt.Sprintf("Level: %d XP: %d", experience.Level, experience.CurrentXP)
	return a.renderText(title, offsetX, offsetY)
}

func (a *renderingSystem) renderMap(em ecs.EntityManager, offsetX, offsetY int, size components.Size) {

	// отрисовываем фон
	for y := 0; y < size.Height; y++ {
		for x := 0; x < size.Width; x++ {
			(*a.screen).SetContent(x+offsetX, y+offsetY, '.', nil, tcell.StyleDefault)
		}
	}

	// отрисовываем сущности имеющие текстуру
	renderable := em.FilterByMask(components.MaskTexture | components.MaskPosition)
	for _, entity := range renderable {
		entityPosition := entity.Get(components.MaskPosition).(*components.Position)
		entityTexture := entity.Get(components.MaskTexture).(*components.Texture)
		(*a.screen).SetContent(entityPosition.X+offsetX, entityPosition.Y+offsetY, rune(*entityTexture), nil, tcell.StyleDefault)
	}
}

func (a *renderingSystem) renderEntityDescription(em ecs.EntityManager, offsetX, offsetY int) (int, int) {
	description := ""
	entities := em.FilterByMask(components.MaskHealth | components.MaskAttack)
	for _, entity := range entities {
		entityID := entity.Id
		entityHealth := entity.Get(components.MaskHealth).(*components.Health)
		entityAttack := entity.Get(components.MaskAttack).(*components.Attack)

		description += fmt.Sprintf("%s:\tHealth: %d\tAttack: %d\n", entityID, entityHealth.CurrentHealth, entityAttack.Damage)
	}
	lines, maxLength := a.renderText(description, offsetX, offsetY)
	return offsetX + maxLength, offsetY + lines
}

// Отрисовка текста
// Возвращает количество строк и длина самой длинной строки
func (a *renderingSystem) renderText(text string, offsetX, offsetY int) (int, int) {
	lines := strings.Split(text, "\n")
	maxLength := 0
	for row, line := range lines {
		maxLength = max(maxLength, len(line))
		for col, char := range line {
			(*a.screen).SetContent(col+offsetX, row+offsetY, char, nil, tcell.StyleDefault)
		}
	}
	return len(lines), maxLength
}

func (a *renderingSystem) Teardown() {}

func (a *renderingSystem) WithScreen(screen *tcell.Screen) *renderingSystem {
	a.screen = screen
	return a
}

func NewRenderingSystem() *renderingSystem {
	return &renderingSystem{}
}
