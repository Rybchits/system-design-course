package systems

import (
	ecs "roguelike/packages/ecs"

	"github.com/gdamore/tcell/v2"
)

type inputSystem struct {
	screen        *tcell.Screen
	eventChan     chan tcell.Event
	inputHandlers []InputHandler
}

type InputHandler interface {
	CanHandle(event *tcell.EventKey) bool

	// возвращает true если нужно ли продолжать игру
	Handle(event *tcell.EventKey, em ecs.EntityManager) bool
}

func (a *inputSystem) Process(em ecs.EntityManager) (engineState int) {
	ev := <-a.eventChan
	switch ev := ev.(type) {
	case *tcell.EventResize:
		(*a.screen).Sync()

	case *tcell.EventKey:
		for _, handler := range a.inputHandlers {
			if handler.CanHandle(ev) {
				engineContinue := handler.Handle(ev, em)
				if !engineContinue {
					return ecs.StateEngineStop
				}
				break
			}
		}
	}
	return ecs.StateEngineContinue
}

func (a *inputSystem) WithScreen(screen *tcell.Screen) *inputSystem {
	a.screen = screen
	return a
}

func (a *inputSystem) WithInputHandlers(handlers ...InputHandler) *inputSystem {
	a.inputHandlers = handlers
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
