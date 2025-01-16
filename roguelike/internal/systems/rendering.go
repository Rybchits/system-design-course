package systems

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"

	"github.com/gdamore/tcell/v2"
)

// Отрисовывает карту и сущности на ней
type renderingSystem struct {
	// Сцена для отрисовки содержимого
	screen *tcell.Screen

	// Текст над полем
	title string

	// Размер игрового поля
	width, height int
}

func (a *renderingSystem) Setup() {}

func (a *renderingSystem) Process(em ecs.EntityManager) (state int) {
	(*a.screen).Clear()

	// отрисовываем фон
	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			(*a.screen).SetContent(x, y, '.', nil, tcell.StyleDefault)
		}
	}

	// отрисовываем сущности имеющие текстуру
	renderable := em.FilterByMask(components.MaskTexture | components.MaskPosition)
	for _, entity := range renderable {
		entityPosition := entity.Get(components.MaskPosition).(*components.Position)
		entityTexture := entity.Get(components.MaskTexture).(*components.Texture)
		(*a.screen).SetContent(entityPosition.X, entityPosition.Y, rune(*entityTexture), nil, tcell.StyleDefault)
	}

	(*a.screen).Show()
	return ecs.StateEngineContinue
}

func (a *renderingSystem) Teardown() {}

func (a *renderingSystem) WithScreen(screen *tcell.Screen) *renderingSystem {
	a.screen = screen
	return a
}

func (a *renderingSystem) WithTitle(title string) *renderingSystem {
	a.title = title
	return a
}

func (a *renderingSystem) WithHeight(height int) *renderingSystem {
	a.height = height
	return a
}

func (a *renderingSystem) WithWidth(width int) *renderingSystem {
	a.width = width
	return a
}

func NewRenderingSystem() *renderingSystem {
	return &renderingSystem{}
}
