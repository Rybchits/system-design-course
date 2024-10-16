package command_meta

import (
	"fmt"
	envsholder "shell/internal/envs_holder"
)

// Структура, хранящая вспомогательную информацию о команде.
// Данная информация используется конкретными структурами, которые реализуют
// интерфейс Command, чтобы исполнить команду.
type CommandMeta struct {
	// Имя команды
	Name string
	// Аргументы команды
	Args []string
	// Локальные для команды переменные окружения
	Envs envsholder.Env
}

func (m *CommandMeta) IsEmpty() bool {
	return m.Name == ""
}

func (m *CommandMeta) Equal(r *CommandMeta) bool {
	if m.Name != r.Name {
		return false
	}

	if len(m.Args) != len(r.Args) {
		return false
	}

	for i := range r.Args {
		if m.Args[i] != r.Args[i] {
			return false
		}
	}

	fmt.Println(m.Envs.Vars)
	fmt.Println(r.Envs.Vars)
	if len(m.Envs.Vars) != len(r.Envs.Vars) {
		return false
	}

	for i := range r.Envs.Vars {
		if m.Envs.Vars[i] != r.Envs.Vars[i] {
			return false
		}
	}

	return true
}
