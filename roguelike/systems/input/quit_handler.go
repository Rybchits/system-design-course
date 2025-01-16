package systems

import (
	ecs "roguelike/esc"

	"github.com/gdamore/tcell/v2"
)

type QuitHandler struct{}

func (h *QuitHandler) CanHandle(event *tcell.EventKey) bool {
	return event.Rune() == 'q'
}

func (h *QuitHandler) Handle(event *tcell.EventKey, em ecs.EntityManager) bool {
	return false
}

func NewQuitHandler() *QuitHandler {
	return &QuitHandler{}
}
