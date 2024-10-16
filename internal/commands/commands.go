package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"shell/internal/command_meta"
	"strings"
)

// Интерфейс, который реализуют все команды,
// представленные в интерпретаторе
type Command interface {
	Execute() error
}

//////////////////////////////////

// Фабрика для создания конкретных команд на основании метаданных команды
type CommandFactory struct {
}

// Метод фабрики, который создает конкретную команду на основании метаданных
func (f *CommandFactory) CommandFromMeta(meta command_meta.CommandMeta, in *os.File, out *os.File) Command {
	switch meta.Name {
	case "cat":
		return CatCommand{in, out, meta}
	case "wc":
		return WcCommand{in, out, meta}
	case "echo":
		return EchoCommand{in, out, meta}
	case "pwd":
		return PwdCommand{in, out, meta}
	case "exit":
		return ExitCommand{in, out, meta}
	default:
		return ProcessCommand{in, out, meta}
	}
}

//////////////////////////////////

// Команда wc.
// Дескрипторами файлов данная структура не владеет.
type WcCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Команда wc выводит количество строк, слов и байтов в файле.
// Имя файла берется из метаданных команды.
// Результат работы выводится в файл, который представлен дескриптором output.
func (cmd WcCommand) Execute() error {
	in := cmd.input
	filename := ""
	if len(cmd.meta.Args) != 0 {
		filename = cmd.meta.Args[0]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("wc: Failed to open file with err: %s\n", err)
			return err
		}
		in = file
		defer file.Close()
	}

	scanner := bufio.NewScanner(in)
	lineCount, wordCount, byteCount := 0, 0, 0

	for scanner.Scan() {
		lineCount++
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
		byteCount += len(scanner.Text()) + 1
	}

	buffer := []byte{}
	if len(filename) != 0 {
		buffer = []byte(fmt.Sprintf("\t%d\t%d\t%d\t%s\n", lineCount, wordCount, byteCount, filename))
	} else {
		buffer = []byte(fmt.Sprintf("\t%d\t%d\t%d\n", lineCount, wordCount, byteCount))
	}

	if _, err := cmd.output.Write(buffer); err != nil {
		fmt.Printf("wc: Failed to write in file with err: %s\n", err)
		return err
	}
	return nil
}

//////////////////////////////////

// Команда cat.
// Дескрипторами файлов данная структура не владеет.
type CatCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Команда cat выводит содержимое файла.
// Имя файла берется из метаданных команды.
// Результат работы выводится в файл, который представлен дескриптором output.
func (cmd CatCommand) Execute() error {
	in := cmd.input
	var err error = nil

	if len(cmd.meta.Args) != 0 {
		filename := cmd.meta.Args[0]
		in, err = os.Open(filename)
	} else {
		in = cmd.input
	}

	if err != nil {
		fmt.Printf("cat: Failed to open file with err: %s\n", err)
		return err
	}

	buffer := make([]byte, 4096)
	for err == nil {
		var n int
		n, err = in.Read(buffer)
		if err == nil {
			_, err = cmd.output.Write(buffer[:n])
		}
	}
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		fmt.Printf("cat: Failed to write in file with err: %s\n", err)
		return err
	}
	return nil
}

//////////////////////////////////

// Команда echo.
// Дескрипторами файлов данная структура не владеет.
type EchoCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Команда echo выводит свои аргументы.
// Результат работы выводится в файл, который представлен дескриптором output.
// Аргументы команды берутся из метаданных команды.
func (cmd EchoCommand) Execute() error {
	buffer := []byte(strings.Join(cmd.meta.Args, " "))
	buffer = append(buffer, '\n')
	if _, err := cmd.output.Write(buffer); err != nil {
		fmt.Printf("echo: Failed to write in file with err: %s\n", err)
		return err
	}
	return nil
}

//////////////////////////////////

// Команда pwd.
// Дескрипторами файлов данная структура не владеет.
type PwdCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Команда pwd выводит содержимое текущей директории.
// Имя директории берется из метаданных команды.
// Результат работы выводится в файл, который представлен дескриптором output.
func (cmd PwdCommand) Execute() error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: Failed to read current directory path with err: %s\n", err)
		return err
	}

	buffer := []byte(dir)
	if _, err := cmd.output.Write(buffer); err != nil {
		fmt.Printf("pwd: Failed to write in file with err: %s\n", err)
		return err
	}
	return nil
}

//////////////////////////////////

// Команда process.
// Дескрипторами файлов данная структура не владеет.
type ProcessCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Данный метод запускает внешнюю программу с указанным именем и набором аргументов.
// Аргументы и имя программы берется из метаданных команды.
// Ввод команда берет из файла, который представлен дескриптором input.
// Результат работы выводится в файл, который представлен дескриптором output.
func (cmd ProcessCommand) Execute() error {
	process := exec.Command(cmd.meta.Name, cmd.meta.Args...)
	process.Stdin = cmd.input
	process.Stdout = cmd.output
	err := process.Run()
	if err != nil {
		fmt.Printf("process: Failed to process command with err: %s\n", err)
		return err
	}
	return nil
}

//////////////////////////////////

// Команда exit.
// Дескрипторами файлов данная структура не владеет.
type ExitCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

// Команда exit завершает исполнение процесса shell.
func (cmd ExitCommand) Execute() error {
	os.Exit(0)
	return nil
}
