package commands

import (
	"fmt"
	"io/fs"
	"os"
	"shell/internal/command_meta"
)

// ListDirCommand выводит файлы и директории в директории-аргументе
// если переданной директории не существует, то будет возвращена ошибка.
// Дескрипторами файлов данная структура не владеет.
type ListDirCommand struct {
	output *os.File
	meta   command_meta.CommandMeta
}

type listDirOptions struct {
	Positional struct {
		Path string
	} `positional-args:"true" maximum:"1"`
}

var _ Command = ListDirCommand{}

// Execute implements Command.
func (cmd ListDirCommand) Execute() error {
	var opts listDirOptions
	err := arg_parse(&opts, cmd.meta.Args)
	if err != nil {
		return err
	}

	// Если нам не передали path, то используем текущую директорию
	path := opts.Positional.Path
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to determine current directory: %v", err)

		}
	}

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("unable to read file info (%s): %v", path, err)
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsRegular():
		if _, err := cmd.output.WriteString(reportEntry(fileInfo)); err != nil {
			return err
		}
	case mode.IsDir():
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("unable to read given directory (%s): %v", path, err)
		}

		for _, entry := range entries {
			fileInfo, err := entry.Info()
			if err != nil {
				return fmt.Errorf("unable to read file info (%s): %v", path, err)
			}

			if _, err := cmd.output.WriteString(reportEntry(fileInfo)); err != nil {
				return err
			}
		}
	}

	return nil
}

func reportEntry(fileInfo fs.FileInfo) string {
	return fmt.Sprintf("%s %s\n", permissionString(fileInfo.Mode()), fileInfo.Name())
}

// permissionString generates an ls-like permission string
func permissionString(mode os.FileMode) string {
	perm := []byte{'-', '-', '-', '-', '-', '-', '-', '-', '-', '-'}

	if mode.IsDir() {
		perm[0] = 'd'
	}
	if mode&os.ModeSymlink != 0 {
		perm[0] = 'l'
	}

	roles := []struct {
		read  int
		write int
		exec  int
	}{
		{1, 2, 3}, // Owner
		{4, 5, 6}, // Group
		{7, 8, 9}, // Others
	}

	// Loop through the roles (owner, group, others)
	for i, role := range roles {
		
		offset := 8 - i * 3 // смещение старшего бита группы
		if mode&(1<<(offset)) != 0 {
			perm[role.read] = 'r'
		}
		if mode&(1<<(offset-1)) != 0 {
			perm[role.write] = 'w'
		}
		if mode&(1<<(offset-2)) != 0 {
			perm[role.exec] = 'x'
		}
	}

	return string(perm)
}
