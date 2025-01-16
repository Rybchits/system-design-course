package game

// Строит игру на основанаии уровня
type GameBuilder interface {
	SetLocation(location string)
	BuildScreen()
	BuildEngine()
	GetResult() GameModel
}
