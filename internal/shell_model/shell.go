package shellmodel

import (
	"fmt"
	"os"
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
			fmt.Println("Commands: ", commands)
			pipeline := self.pipelineFactory.CreatePipeline(input, output, commands)
			err = pipeline.Execute()
			if err != nil {
				fmt.Printf("Issue running pipeline %s\n", err)
			}
		}
	}
}

func (self *Shell) Terminate() {
	self.terminate <- true
}
