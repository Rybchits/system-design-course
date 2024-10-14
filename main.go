package main

import (
	"fmt"
	"os"
	"shell/internal/parser"
)

func main() {
	/*l := parser.NewTokenizer(os.Stdin)
	token, err := l.Next()

	for err == nil {
		println("New token: ", token.Value, " ", token.TokenType)
		token, err = l.Next()
	}
	println("error ", err)*/

	a := parser.NewParser(os.Stdin)
	a.Init()

	for event := range a.Listen() {
		if event.Error != "" {
			fmt.Printf("Error: %s\n", event.Error)
		} else {
			fmt.Printf("Name: %s\n", event.Command.Name)
			for _, arg := range event.Command.Args {
				fmt.Printf("Arg: %s ", arg)
			}
		}
	}
	a.Dispose()
}
