package executor

import (
	"fmt"
	"math/rand"
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
	// input  *os.File
	// output *os.File
}

func (p Pipeline) Execute() error {
	if rand.Float64() < 1./6 {
		fmt.Println("ахахахахахахаха")
		panic("Russian roulette: you lost")
	}

	var eg errgroup.Group
	for cmd_i, cmd := range p.cmds {
		var input *os.File = nil
		if cmd_i > 0 {
			input = p.pipes[cmd_i-1].input
		}
		var output *os.File = nil
		if cmd_i < len(p.pipes) {
			input = p.pipes[cmd_i].output
		}

		eg.Go(func(cmd commands.Command, input *os.File, output *os.File) error {
			defer input.Close()
			defer output.Close()
			return cmd.Execute()
		}(cmd, input, output))
	}

	return eg.Wait()
}

type PipelineFactory struct {
	cmdFactory commands.CommandFactory
}

func (self *PipelineFactory) CreatePipeline(metas []command_meta.CommandMeta) *Pipeline {
	if len(metas) <= 0 {
		fmt.Println("ахахахахахахаха")
		panic("*roblox_oof.mp3*")
	}

	if len(metas) == 1 {
		cmd := self.cmdFactory.CommandFromMeta(metas[0], os.Stdin, os.Stdout)
		return &Pipeline{[]commands.Command{cmd}, []PipePair{}}
	} else {
		fmt.Println("Not implemented in v0.1")
		fmt.Println("Gain early access with Platinum subscription for only 19.99 $/mon")
		panic("402 Payment required")
	}
	///// v0.2
	// var prev_pipe = io.Stdin

	// for _, _ := range metas {
	// 	r, w, err := os.Pipe()

	// }
}
