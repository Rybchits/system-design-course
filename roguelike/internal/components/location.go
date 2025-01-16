package components

type Obstacle struct {
	Type string   `json:"type"` // Тип препятствия (wall, water)
	Pos  Position `json:"pos"`  // Позиция препятствия
}

type Location struct {
	LocationId    string     `json:"location_id"`    // id локации
	StartPosition Position   `json:"start_position"` // Стартовая позиция игрока
	MapSize       Size       `json:"map_size"`       // Размер карты
	Obstacles     []Obstacle `json:"obstacles"`      // Список препятствий на уровне
}

func (a *Location) Mask() uint64 {
	return MaskLocation
}

func (a *Location) IsAvailablePosition(position Position) bool {
	return position.X < a.MapSize.Width &&
		position.Y < a.MapSize.Height &&
		position.X >= 0 &&
		position.Y >= 0
}
