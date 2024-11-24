package commands

import (
	"os"
	"shell/internal/command_meta"
)

// ChangeDirCommand изменяет текущую рабочую директорию терминала.
// Если команда вызвана без аргументов, то текущей рабочей директорией становится домашнаяя директория пользователя
// Если у пользователя не установленая домашная директория - команда вернет ошибку.
type ChangeDirCommand struct {
	meta command_meta.CommandMeta
}

type сhangeDirOptions struct {
	Positional struct {
		Path string
	} `positional-args:"true" maximum:"1"`
}

var _ Command = ChangeDirCommand{}

func (cmd ChangeDirCommand) Execute() error {
	var opts сhangeDirOptions
	err := arg_parse(&opts, cmd.meta.Args)
	if err != nil {
		return err
	}

	path := opts.Positional.Path
	if path == "" {
		path, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}

	return os.Chdir(path)
}
