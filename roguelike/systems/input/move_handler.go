package systems

import (
	"roguelike/components"
	ecs "roguelike/esc"

	"github.com/gdamore/tcell/v2"
)

type MoveHandler struct{}

func (h *MoveHandler) CanHandle(event *tcell.EventKey) bool {
	available := []rune{'w', 'W', 'a', 'A', 's', 'S', 'd', 'D'}
	symbol := event.Rune()
	for _, v := range available {
		if v == symbol {
			return true
		}
	}
	return false
}

func canMove(newPosition *components.Position, width, height int) bool {
	return newPosition.X < width &&
		newPosition.Y < height &&
		newPosition.X >= 0 &&
		newPosition.Y >= 0
}

func (h *MoveHandler) Handle(event *tcell.EventKey, em ecs.EntityManager) bool {
	player := em.Get("player")
	playerPosition := player.Get(components.MaskPosition).(*components.Position)

	location := em.Get("location").Get(components.MaskLocation).(*components.Location)
	newPlayerPosition := playerPosition.Clone()

	switch event.Rune() {
	case 'w', 'W':
		newPlayerPosition.WithY(playerPosition.Y - 1)
	case 'a', 'A':
		newPlayerPosition.WithX(playerPosition.X - 1)
	case 's', 'S':
		newPlayerPosition.WithY(playerPosition.Y + 1)
	case 'd', 'D':
		newPlayerPosition.WithX(playerPosition.X + 1)
	}
	if canMove(newPlayerPosition, location.MapSize.Width, location.MapSize.Height) {
		*playerPosition = *newPlayerPosition
	}
	return true
}

func NewMoveHandler() *MoveHandler {
	return &MoveHandler{}
}
