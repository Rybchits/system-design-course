package shellmodel

import (
	"fmt"
	"io"
	"os"
	envsholder "shell/internal/envs_holder"
	"shell/internal/executor"
	"shell/internal/parser"
)

type Shell struct {
	pipelineFactory *executor.PipelineFactory
	terminate       chan bool
}

func NewShell() *Shell {
	pipelineFactory := executor.NewPipelineFactory()
	return &Shell{pipelineFactory: pipelineFactory, terminate: make(chan bool)}
}

// Основной цикл оболочки
// Обрабатывает пользовательский ввод
func (self *Shell) ShellLoop(input *os.File, output *os.File, to_greet bool) {

	tokenizer := parser.NewTokenizer(input, &envsholder.GlobalEnv)
	curr_parser := parser.NewParser(tokenizer)
	for {
		if to_greet {
			output.WriteString("$ ")
		}
		commands, err := curr_parser.Parse()
		end_of_file := err == io.EOF

		if err != nil && !end_of_file {
			output.WriteString("Parse issue\n")
			continue
		}

		pipeline := self.pipelineFactory.CreatePipeline(input, output, commands)
		if pipeline != nil {
			err = pipeline.Execute()
			if err != nil {
				envsholder.GlobalEnv.Set(envsholder.ExecStatusKey, "1")
				fmt.Printf("%s\n", err)
			}
		}

		if end_of_file {
			return
		}
	}
}

func (self *Shell) Terminate() {
	os.Exit(0)
}
