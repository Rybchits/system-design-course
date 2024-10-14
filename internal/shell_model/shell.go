package shellmodel

import (
	"fmt"
	"os"
	"shell/internal/parser"
)

// type State int;

// const (
// 	INIT State = iota
// 	EXECUTING
// 	EXIT
// )

// type Shell struct {
// 	state State
// }

// func NewShell() *Shell {
// 	return &Shell { INIT }
// }

func //(self *Shell)
ShellLoop() {
	curr_parser := parser.NewParser(os.Stdin)

	//self.state = EXECUTING

	for {
		fmt.Print("Enter command: ")
		command := curr_parser.ParseCommand() //?
		/*
			command -> executor factory ?
			factory.GetPipeline()
			res, err := pipeline.Execute()
			if err != Nil {
				fmt.Println("Pizdos")
			} else {
				fmt.Println("zaebis")
			}
		*/
		// command -> executor factory ?
		// factory.GetPipeline()
		// pipeline
	}

}
