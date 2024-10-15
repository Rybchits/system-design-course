package envsholder

// Ассоциативный контейнер - хранилище переменных окружения
type Env struct {
	vars map[string]string
}
