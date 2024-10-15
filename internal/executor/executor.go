package executor

import (
	"os"
	"shell/internal/command_meta"
	"shell/internal/commands"

	"golang.org/x/sync/errgroup"
)

type PipePair struct {
	input  *os.File
	output *os.File
}

type Pipeline struct {
	cmds  []commands.Command
	pipes []PipePair
}

func (p Pipeline) Execute() error {
	var eg errgroup.Group
	for cmd_i, cmd := range p.cmds {
		var input *os.File = nil
		if cmd_i > 0 {
			input = p.pipes[cmd_i-1].input
		}
		var output *os.File = nil
		if cmd_i < len(p.pipes) {
			output = p.pipes[cmd_i].output
		}

		cmddd := cmd
		eg.Go(func() error {
			defer func() {
				if input != nil {
					input.Close()
				}
			}()
			defer func() {
				if output != nil {
					output.Close()
				}
			}()
			return cmddd.Execute()
		})
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

// Создает пайплайн исполнения на основе переданных информации о командах
func (self *PipelineFactory) CreatePipeline(input *os.File, output *os.File, metas []command_meta.CommandMeta) *Pipeline {
	if len(metas) <= 0 {
		return nil
		panic("*roblox_oof.mp3*")
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
