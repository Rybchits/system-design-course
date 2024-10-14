package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"shell/internal/command_meta"
	"strings"
)

type Command interface {
	Execute() error
}

//////////////////////////////////

type CommandFactory struct {
}

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

type WcCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

func (cmd WcCommand) Execute() error {
	if len(cmd.meta.Args) != 1 {
		fmt.Printf("wc: Wrong count of arguments: %s\n", os.ErrInvalid)
		return os.ErrInvalid
	}

	filename := cmd.meta.Args[0]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("wc: Failed to open file with err: %s\n", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount, wordCount, byteCount := 0, 0, 0

	for scanner.Scan() {
		lineCount++
		words := strings.Fields(scanner.Text())
		wordCount += len(words)
		byteCount += len(scanner.Text())
	}

	buffer := []byte(fmt.Sprintf("%d %d %d %s\n", lineCount, wordCount, byteCount, filename))
	cmd.output.Write(buffer)
	return nil
}

//////////////////////////////////

type CatCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

func (cmd CatCommand) Execute() error {
	if len(cmd.meta.Args) != 1 {
		fmt.Printf("cat: Wrong count of arguments: %s\n", os.ErrInvalid)
		return os.ErrInvalid
	}

	filename := cmd.meta.Args[0]
	buffer, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("cat: Failed to open file with err: %s\n", err)
		return err
	}

	cmd.output.Write(buffer)
	return nil
}

//////////////////////////////////

type EchoCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

func (cmd EchoCommand) Execute() error {
	buffer := []byte(strings.Join(cmd.meta.Args, " "))
	cmd.output.Write(buffer)
	return nil
}

//////////////////////////////////

type PwdCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

func (cmd PwdCommand) Execute() error {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: Failed to read directory with err: %s\n", err)
		return err
	}

	buffer := []byte(dir)
	cmd.output.Write(buffer)
	return nil
}

//////////////////////////////////

type ProcessCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

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

type ExitCommand struct {
	input  *os.File
	output *os.File
	meta   command_meta.CommandMeta
}

func (cmd ExitCommand) Execute() error {
	os.Exit(0)
	return nil
}
