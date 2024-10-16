package command_meta

import (
	"maps"
	"slices"
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
	return m.Name == r.Name && maps.Equal(m.Envs.Vars, r.Envs.Vars) && slices.Equal(m.Args, r.Args);
}
