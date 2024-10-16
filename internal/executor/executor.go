package executor

import (
	"os"
	"shell/internal/command_meta"
	"shell/internal/commands"

	"golang.org/x/sync/errgroup"
)

// Пара дескрипторов на чтение и на запись
type PipePair struct {
	input  *os.File
	output *os.File
}

// Набор команд, соединенных пайпами
type Pipeline struct {
	cmds  []commands.Command
	pipes []PipePair
}

// Выполнить пайплайн из команд.
// В случае ошибки какой-либо из команд пайплайна, все остальные завершают свою работу.
// Значение ошибки возвращается в вызывающую функцию.
func (p Pipeline) Execute() error {
	var eg errgroup.Group
	chs := make([]chan bool, len(p.cmds))
	for ch_i := range chs {
		chs[ch_i] = make(chan bool)
	}

	for cmd_i, cmd := range p.cmds {
		eg.Go(func() error {
			defer func() {
				chs[cmd_i] <- true
			}()
			res := cmd.Execute()
			return res
		})
	}

	for cmd_i := range p.cmds {
		<-chs[cmd_i]
		if cmd_i > 0 {
			p.pipes[cmd_i-1].input.Close()
		}
		if cmd_i < len(p.pipes) {
			p.pipes[cmd_i].output.Close()
		}
	}

	return eg.Wait()
}

type PipelineFactory struct {
	cmdFactory *commands.CommandFactory
}

func NewPipelineFactory() *PipelineFactory {
	cmFactory := commands.CommandFactory{}
	return &PipelineFactory{cmdFactory: &cmFactory}
}

// Создает пайплайн исполнения на основе переданной информации о командах
func (self *PipelineFactory) CreatePipeline(input *os.File, output *os.File, metas []command_meta.CommandMeta) *Pipeline {
	if len(metas) <= 0 {
		return nil
	}

	var pipeline *Pipeline = &Pipeline{}
	fokgobak := false
	for i := 0; i < len(metas); i++ {
		if i < len(metas)-1 {
			r, w, err := os.Pipe()
			if err != nil {
				fokgobak = true
				break
			}
			pipeline.pipes = append(pipeline.pipes, PipePair{r, w})
		}
		in := input
		if i > 0 {
			in = pipeline.pipes[i-1].input
		}
		out := output
		if i < len(metas)-1 {
			out = pipeline.pipes[i].output
		}
		cmd := self.cmdFactory.CommandFromMeta(metas[i], in, out)
		pipeline.cmds = append(pipeline.cmds, cmd)
	}

	if fokgobak {
		for _, pipe := range pipeline.pipes {
			pipe.input.Close()
			pipe.output.Close()
		}
		pipeline = nil
	}

	return pipeline
}
