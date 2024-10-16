package parser

import (
	"io"
	"strings"
	"shell/internal/command_meta"
)

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(read io.Reader) *Parser {
	return &Parser{
		tokenizer: NewTokenizer(read),
	}
}

func (p *Parser) Parse() ([]command_meta.CommandMeta, error) {
	pipe := make([]command_meta.CommandMeta, 0)
	current := command_meta.CommandMeta{}
	for {
		token, err := p.tokenizer.Next()
		if err == nil {
			switch token.TokenType {
			case WordToken:
				{
					if current.Name == "" {
						if (strings.Contains(token.Value, "=")) {
							parts := strings.SplitN(token.Value, "=", 2)
							current.Envs.Vars[parts[0]] = parts[1]
						} else {
							current.Name = token.Value
						}
					} else {
						current.Args = append(current.Args, token.Value)
					}
				}
			case PipeToken:
				{
					if !current.IsEmpty() {
						pipe = append(pipe, current)
						current = command_meta.CommandMeta{}
					}
				}
			case EndLineToken:
				{
					if !current.IsEmpty() {
						pipe = append(pipe, current)
					}

					if len(pipe) != 0 {
						return pipe, nil
					}
				}
			}
		} else if err == io.EOF {
			return pipe, nil
		} else {
			return pipe, err
		}
	}
}
