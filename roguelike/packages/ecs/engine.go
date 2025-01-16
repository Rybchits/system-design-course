package ecs

// состояние Engine
const (
	StateEngineContinue = 0
	StateEngineStop     = 1
)

// обрабатывает этапы Setup(), Run() и Teardown() для всех систем
type Engine interface {
	// Запускает метод Process для всех System
	// пока ShouldEngineStop не будет установлен в true
	Run()

	// Запускает метод Setup для всех System
	// и инициализирует ShouldEngineStop и ShouldEnginePause с false
	Setup()

	// Teardown вызывает метод Teardown() для каждой системы
	Teardown()

	// Tick вызывает метод Process() для каждой системы ровно один раз
	Tick()
}
