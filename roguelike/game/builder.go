package game

// Строит игру на основанаии уровня
type GameBuilder interface {
	SetLevel(level int)
	BuildScreen()
	BuildEngine()
	GetResult() GameModel
}
