package parser

import (
	"fmt"
	"io"
	"shell/internal/command_meta"
)

type Parser struct {
	tokenizer *Tokenizer
	channel   chan ParserEvent
}

type ParserEvent struct {
	Command command_meta.CommandMeta
	Error   string
}

func NewParser(read io.Reader) *Parser {
	tokenizer := NewTokenizer(read)
	channel := make(chan ParserEvent)
	return &Parser{tokenizer: tokenizer, channel: channel}
}

func (p *Parser) Listen() <-chan ParserEvent {
	return p.channel
}

func mapFunc[T any, R any](input []T, f func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

func (p *Parser) Init() {
	tokens := []Token{}
	for {
		token, err := p.tokenizer.Next()
		fmt.Println("dub ", token)
		if err != nil {
			fmt.Println("err ", err.Error())
			if err == io.EOF {
				if len(tokens) == 0 {
					continue
				}
				name := tokens[0].Value
				args := mapFunc(tokens[1:], func(token Token) string { return token.Value })
				p.channel <- ParserEvent{Command: command_meta.CommandMeta{Name: name, Args: args}}
				tokens = []Token{}
			} else {
				p.channel <- ParserEvent{Error: err.Error()}
			}
		} else if token.TokenType == WordToken {
			tokens = append(tokens, *token)

		} else if token.TokenType == SpaceToken {
			fmt.Println("Space ", token.Value)
		}
	}
}

func (p *Parser) Dispose() {
	close(p.channel)
}
