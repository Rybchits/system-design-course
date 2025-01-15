package systems

import (
	"roguelike/components"

	ecs "roguelike/esc"

	"github.com/gdamore/tcell/v2"
)

type inputSystem struct {
	screen    *tcell.Screen
	eventChan chan tcell.Event
}

func canMove(newPosition *components.Position, width, height int) bool {
	return newPosition.X < width && newPosition.Y < height && newPosition.X >= 0 && newPosition.Y >= 0
}

func (a *inputSystem) Process(em ecs.EntityManager) (engineState int) {
	ev := <-a.eventChan
	switch ev := ev.(type) {
	case *tcell.EventResize:
		(*a.screen).Sync()

	case *tcell.EventKey:
		player := em.Get("player")

		playerPosition := player.Get(components.MaskPosition).(*components.Position)
		newPlayerPosition := playerPosition.Clone()

		switch ev.Rune() {
		case 'w', 'W':
			newPlayerPosition.WithY(playerPosition.Y - 1)
		case 'a', 'A':
			newPlayerPosition.WithX(playerPosition.X - 1)
		case 's', 'S':
			newPlayerPosition.WithY(playerPosition.Y + 1)
		case 'd', 'D':
			newPlayerPosition.WithX(playerPosition.X + 1)
		case 'q':
			return ecs.StateEngineStop
		}
		// TODO Вынести логику перехода
		if canMove(newPlayerPosition, 10, 10) {
			*playerPosition = *newPlayerPosition
		}
	}
	return ecs.StateEngineContinue
}

func (a *inputSystem) WithScreen(screen *tcell.Screen) *inputSystem {
	a.screen = screen
	return a
}

func (a *inputSystem) Setup() {
	a.eventChan = make(chan tcell.Event)
	go func() {
		for {
			ev := (*a.screen).PollEvent()
			if ev != nil {
				a.eventChan <- ev
			}
		}
	}()
}

func (a *inputSystem) Teardown() {
	close(a.eventChan)
}

func NewInputSystem() *inputSystem {
	return &inputSystem{}
}
