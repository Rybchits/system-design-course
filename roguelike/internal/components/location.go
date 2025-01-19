package components

// Конфигурация противника
type Enemy struct {
	Type   string   `json:"type"`   // Тип врага (aggressive, passive)
	Pos    Position `json:"pos"`    // Позиция врага
	Health int      `json:"health"` // Здоровье врага
	Attack int      `json:"attack"` // Атака врага
}

// Компонент локации, описывает стартовые параметры локации
type Location struct {
	LocationId    string   `json:"location_id"`    // id локации
	StartPosition Position `json:"start_position"` // Стартовая позиция игрока
	MapSize       Size     `json:"map_size"`       // Размер карты
	Enemies       []Enemy  `json:"enemies"`        // Список врагов на уровне
}

func (a *Location) Mask() uint64 {
	return MaskLocation
}

// Проверяет является ли позиция допустимой для текущей локации
func (a *Location) IsAvailablePosition(position Position) bool {
	return position.X < a.MapSize.Width &&
		position.Y < a.MapSize.Height &&
		position.X >= 0 &&
		position.Y >= 0
}
