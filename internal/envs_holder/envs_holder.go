package envsholder

import (
	"fmt"
)

// Ассоциативный контейнер - хранилище переменных окружения
type Env struct {
	Vars map[string]string
}

// Получить все переменные окружения в виде набора строк вида "ключ=значение"
func (e *Env) Environ() []string {
	result := make([]string, 0, len(e.Vars))
	for key, value := range e.Vars {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}
	return result
}

// Установить значение переменной окружения
func (e *Env) Set(key string, value string) {
	e.Vars[key] = value
}

// Очистить все переменные окружения
func (e *Env) Clear() {
	e.Vars = make(map[string]string)
}

//////////////////////////////////

const (
	ExecStatusKey = "?"
	OkStatusValue = "0"
)

//////////////////////////////////

// Хранилище переменных окружения
var GlobalEnv = Env{
	Vars:  map[string]string{ExecStatusKey: OkStatusValue},
}
