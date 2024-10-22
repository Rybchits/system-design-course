package parser

import (
	"io"
	"shell/internal/command_meta"
	"strings"
)

type Parser struct {
	tokenizer *Tokenizer
}

func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer: tokenizer,
	}
}

func (p *Parser) Parse() ([]command_meta.CommandMeta, error) {
	pipe := make([]command_meta.CommandMeta, 0)
	current := command_meta.CommandMeta{}
	//var prev_token TokenType = nil

	for {
		token, err := p.tokenizer.Next()
		if err == nil {
			switch token.TokenType {
			case WordToken:
				{
					if current.Name == "" {
						if strings.Contains(token.Value, "=") {
							current.Envs.Init()
							parts := strings.SplitN(token.Value, "=", 2)
							current.Envs.Set(parts[0], parts[1])
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
					return pipe, nil
				}
			}
		} else if err == io.EOF {
			return pipe, io.EOF
		} else {
			return pipe, err
		}
	}
}
