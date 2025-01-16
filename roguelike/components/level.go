package components

type Obstacle struct {
	Type string   `json:"type"` // Тип препятствия (wall, water)
	Pos  Position `json:"pos"`  // Позиция препятствия
}

type Level struct {
	LevelNumber   int        `json:"level_number"`   // Номер текущего уровня
	StartPosition Position   `json:"start_position"` // Стартовая позиция игрока
	MapSize       Size       `json:"map_size"`       // Размер карты
	NextLevelXP   float64    `json:"next_level_xp"`  // Требуемый опыт для следующего уровня
	Obstacles     []Obstacle `json:"obstacles"`      // Список препятствий на уровне
}
