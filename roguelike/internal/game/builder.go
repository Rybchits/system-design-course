package game

// Строит игру на основанаии уровня
type GameBuilder interface {
	SetLocation(locationFilePath string) error
	BuildScreen()
	BuildEngine()
	GetResult() GameModel
}
