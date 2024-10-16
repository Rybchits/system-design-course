package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TokenType int

type runeTokenClass int

type lexerState int

type Token struct {
	TokenType TokenType
	Value     string
}

func (a *Token) Equal(b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.TokenType != b.TokenType {
		return false
	}
	return a.Value == b.Value
}

const (
	spaceRunes            = " \t\r"
	escapingQuoteRunes    = `"`
	nonEscapingQuoteRunes = "'"
	escapeRunes           = `\`
	commentRunes          = "#"
	endLineRunes          = "\n"
	pipeRunes             = "|"
)

const (
	unknownRuneClass runeTokenClass = iota
	spaceRuneClass
	escapingQuoteRuneClass
	nonEscapingQuoteRuneClass
	escapeRuneClass
	commentRuneClass
	endLineRuneClass
	pipeRuneClass
	eofRuneClass
)

const (
	UnknownToken TokenType = iota
	WordToken
	SpaceToken
	CommentToken
	EndLineToken
	PipeToken
)

const (
	startState           lexerState = iota // еще не было символов
	inWordState                            // в процессе определения слова
	escapingState                          // экранирование, следующий символ должен быть литеральным
	escapingQuotedState                    // экранирование в заключенной в кавычки строке
	quotingEscapingState                   // внутри заключенной в кавычки строки, которая поддерживает экранирование
	quotingState                           // внутри строки, которая не поддерживает экранирование
	commentState                           // в пределах комментария
	pipeSymbolState                        // прошлый символ был pipe
	endLineState                           // прошлый символ был \n
)

type tokenClassifier map[rune]runeTokenClass

func (typeMap tokenClassifier) addRuneClass(runes string, tokenType runeTokenClass) {
	for _, runeChar := range runes {
		typeMap[runeChar] = tokenType
	}
}

func newDefaultClassifier() tokenClassifier {
	t := tokenClassifier{}
	t.addRuneClass(spaceRunes, spaceRuneClass)
	t.addRuneClass(escapingQuoteRunes, escapingQuoteRuneClass)
	t.addRuneClass(nonEscapingQuoteRunes, nonEscapingQuoteRuneClass)
	t.addRuneClass(escapeRunes, escapeRuneClass)
	t.addRuneClass(commentRunes, commentRuneClass)
	t.addRuneClass(endLineRunes, endLineRuneClass)
	t.addRuneClass(pipeRunes, pipeRuneClass)
	return t
}

func (t tokenClassifier) ClassifyRune(runeVal rune) runeTokenClass {
	return t[runeVal]
}

type Tokenizer struct {
	input      bufio.Reader
	classifier tokenClassifier
	state      lexerState
	isEnded    bool
}

func NewTokenizer(r io.Reader) *Tokenizer {
	input := bufio.NewReader(r)
	classifier := newDefaultClassifier()
	return &Tokenizer{
		input:      *input,
		classifier: classifier,
		state:      startState,
		isEnded:    false,
	}
}

func (t *Tokenizer) scanStream() (*Token, error) {
	var tokenType TokenType
	var value []rune
	var nextRune rune
	var nextRuneType runeTokenClass
	var err error

	if t.isEnded {
		return nil, io.EOF
	}

	for {
		if t.state == pipeSymbolState {
			t.state = startState
			return &Token{TokenType: PipeToken, Value: pipeRunes}, nil

		} else if t.state == endLineState {
			t.state = startState
			return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
		}

		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil

		} else if err != nil {
			return nil, err
		}

		switch t.state {
		case startState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						return nil, io.EOF
					}
				case spaceRuneClass:
					{
					}
				case escapingQuoteRuneClass:
					{
						tokenType = WordToken
						t.state = quotingEscapingState
					}
				case nonEscapingQuoteRuneClass:
					{
						tokenType = WordToken
						t.state = quotingState
					}
				case escapeRuneClass:
					{
						tokenType = WordToken
						t.state = escapingState
					}
				case commentRuneClass:
					{
						tokenType = CommentToken
						t.state = commentState
					}
				case endLineRuneClass:
					{
						token := &Token{
							TokenType: EndLineToken,
							Value:     string(nextRune)}
						return token, nil
					}
				case pipeRuneClass:
					{
						token := &Token{
							TokenType: PipeToken,
							Value:     string(nextRune)}
						return token, nil
					}
				default:
					{
						tokenType = WordToken
						value = append(value, nextRune)
						t.state = inWordState
					}
				}
			}
		case inWordState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						t.state = startState
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case spaceRuneClass:
					{
						t.state = startState
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case escapingQuoteRuneClass:
					{
						t.state = quotingEscapingState
					}
				case nonEscapingQuoteRuneClass:
					{
						t.state = quotingState
					}
				case endLineRuneClass:
					{
						t.state = endLineState
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, nil
					}
				case pipeRuneClass:
					{
						t.state = pipeSymbolState
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, nil
					}
				case escapeRuneClass:
					{
						t.state = escapingState
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case escapingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found after escape character")
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				default:
					{
						t.state = inWordState
						value = append(value, nextRune)
					}
				}
			}
		case escapingQuotedState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found after escape character")
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				default:
					{
						t.state = quotingEscapingState
						value = append(value, nextRune)
					}
				}
			}
		case quotingEscapingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found when expecting closing quote")
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case escapingQuoteRuneClass:
					{
						t.state = inWordState
					}
				case escapeRuneClass:
					{
						t.state = escapingQuotedState
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case quotingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						err = fmt.Errorf("EOF found when expecting closing quote")
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case nonEscapingQuoteRuneClass:
					{
						t.state = inWordState
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case commentState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						t.state = startState
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case endLineRuneClass:
					{
						t.state = startState
						return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
					}
				case spaceRuneClass:
					{
						if nextRune == '\n' {
							t.state = endLineState
							token := &Token{
								TokenType: tokenType,
								Value:     string(value)}
							return token, err
						} else {
							value = append(value, nextRune)
						}
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		default:
			{
				return nil, fmt.Errorf("Unexpected state: %v", t.state)
			}
		}
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	return t.scanStream()
}

func SplitOnTokens(s string) ([]Token, error) {
	tokenizer := NewTokenizer(strings.NewReader(s))
	tokens := make([]Token, 0)
	for {
		token, err := tokenizer.Next()
		if err != nil {
			if err == io.EOF {
				return tokens, nil
			}
			return []Token{}, err
		}
		tokens = append(tokens, *token)
	}
}
