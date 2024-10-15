package command_meta

// Структура, хранящая вспомогательную информацию о команде.
// Данная информация используется конкретными структурами, которые реализуют
// интерфейс Command, чтобы исполнить команду.
type CommandMeta struct {
	Name string
	Args []string
}

func (m *CommandMeta) IsEmpty() bool {
	return m.Name == ""
}
