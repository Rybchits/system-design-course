package systems

import (
	"roguelike/internal/components"
	ecs "roguelike/packages/ecs"

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

// Считывает нажатие WASD и добавляет игроку компонент передвижения на новую позицию
func (h *MoveHandler) Handle(event *tcell.EventKey, em ecs.EntityManager) bool {
	player := em.Get("player")

	playerPosition := player.Get(components.MaskPosition).(*components.Position)
	newPlayerPosition := playerPosition.Clone()

	location := em.Get("location").Get(components.MaskLocation).(*components.Location)

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

	if location.IsAvailablePosition(*newPlayerPosition) {
		player.Add(components.NewMovement().WithNext(*newPlayerPosition))
	}
	return true
}

func NewMoveHandler() *MoveHandler {
	return &MoveHandler{}
}
