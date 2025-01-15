package systems

import (
	"roguelike/components"
	ecs "roguelike/esc"

	"github.com/gdamore/tcell/v2"
)

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
	for y := 0; y < a.height; y++ {
		for x := 0; x < a.width; x++ {
			(*a.screen).SetContent(x, y, '.', nil, tcell.StyleDefault)
		}
	}
	player := em.Get("player")
	playerPosition := player.Get(components.MaskPosition).(*components.Position)
	playerTexture := player.Get(components.MaskTexture).(*components.Texture)
	(*a.screen).SetContent(playerPosition.X, playerPosition.Y, rune(*playerTexture), nil, tcell.StyleDefault)

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
