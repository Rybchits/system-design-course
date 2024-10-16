package shellmodel

import (
	"fmt"
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
func (self *Shell) ShellLoop(input *os.File, output *os.File) {
	curr_parser := parser.NewParser(input)

	for {
		select {
		case <-self.terminate:
			return
		default:
			commands, err := curr_parser.Parse()
			if err != nil {
				output.WriteString("Parse issue\n")
				continue
			}
			pipeline := self.pipelineFactory.CreatePipeline(input, output, commands)
			if pipeline == nil {
				continue
			}
			err = pipeline.Execute()
			if err != nil {
				envsholder.GlobalEnv.Set(envsholder.ExecStatusKey, "1")
				fmt.Printf("Issue running pipeline %s\n", err)
			}
		}
	}
}

func (self *Shell) Terminate() {
	self.terminate <- true
}
